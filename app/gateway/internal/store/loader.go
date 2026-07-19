package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/Aarav-S2005/mini-api-gateway/app/db/models"
	loadbalancer "github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/lb"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/proxy"
	pluginManager "github.com/Aarav-S2005/mini-api-gateway/app/plugin-manager/registry"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrCouldNotLoadFromDB = errors.New("could not load from database")

var jwtPluginName = "jwt-auth"
var middlewareOrder = map[string]int{
	"cors":       0,
	"ip-limit":   1,
	"rate-limit": 2,
}

func getAllDocuments[T any](ctx context.Context, db *mongo.Database, collectionName string) ([]T, error) {
	collection := db.Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("%w (%s): %w", ErrCouldNotLoadFromDB, collectionName, err)
	}
	defer cursor.Close(ctx)

	var docs []T
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("%w (%s): %w", ErrCouldNotLoadFromDB, collectionName, err)
	}

	return docs, nil
}

func LoadSnapshot(db *mongo.Database, transport *http.Transport, lbManager *loadbalancer.LBManager, registry *pluginManager.PluginRegistry) (*Snapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	projects, err := getAllDocuments[models.Project](ctx, db, "projects")
	if err != nil {
		return nil, err
	}
	upstreams, err := getAllDocuments[models.Upstream](ctx, db, "upstreams")
	if err != nil {
		return nil, err
	}
	routes, err := getAllDocuments[models.Route](ctx, db, "routes")
	if err != nil {
		return nil, err
	}

	snapshot := &Snapshot{
		Projects:  make(map[string]*RuntimeProject),
		Upstreams: make(map[bson.ObjectID]map[string]*RuntimeUpstream),
	}

	if err = loadRuntimeUpstreams(upstreams, snapshot, transport, lbManager); err != nil {
		return nil, err
	}
	if err = loadProjects(projects, routes, snapshot, lbManager, registry); err != nil {
		return nil, err
	}

	return snapshot, nil
}

func loadRuntimeUpstreams(upstreams []models.Upstream, snapshot *Snapshot, transport *http.Transport, lbManager *loadbalancer.LBManager) error {
	for _, upstream := range upstreams {

		backends := make([]RuntimeBackend, 0, len(upstream.Backends))

		for _, backend := range upstream.Backends {

			target, err := url.Parse(backend.URL)
			if err != nil {
				return fmt.Errorf("invalid backend url %q: %w", backend.URL, err)
			}

			weight := 1
			if backend.Weight != nil {
				weight = *backend.Weight
			}

			backends = append(backends, RuntimeBackend{
				URL:    target,
				Proxy:  proxy.BuildReverseProxy(upstream.ID, target, transport, lbManager),
				Weight: weight,
			})
		}

		if snapshot.Upstreams[upstream.ProjectID] == nil {
			snapshot.Upstreams[upstream.ProjectID] = make(map[string]*RuntimeUpstream)
		}

		snapshot.Upstreams[upstream.ProjectID][upstream.Name] = &RuntimeUpstream{
			ID:       upstream.ID,
			Strategy: upstream.LoadBalancingStrategy,
			Backends: backends,
		}
	}
	return nil
}

func loadProjects(projects []models.Project, routes []models.Route, snapshot *Snapshot, lbManager *loadbalancer.LBManager, registry *pluginManager.PluginRegistry) error {
	projectRoutes := make(map[bson.ObjectID][]models.Route)
	for _, route := range routes {
		if !route.Enabled {
			continue
		}

		projectRoutes[route.ProjectID] = append(projectRoutes[route.ProjectID], route)
	}

	for _, project := range projects {

		globalMiddlewares, jwtConfig := SplitJWT(project.Middlewares)

		chain, err := BuildMiddlewareChain(globalMiddlewares, registry)
		if err != nil {
			log.Printf("ERROR: project %s has invalid middleware config, excluding from snapshot: %v", project.ID, err)
			continue
		}

		var jwtMW func(http.Handler) http.Handler
		if jwtConfig != nil {
			plugin, ok := registry.Get(jwtConfig.Name)
			if !ok {
				log.Printf("ERROR: project %s references unknown jwt plugin %q, excluding from snapshot", project.ID, jwtConfig.Name)
				continue
			}
			jwtMW, err = plugin.CreateMiddleware(jwtConfig.Config)
			if err != nil {
				log.Printf("ERROR: project %s jwt middleware config invalid, excluding from snapshot: %v", project.ID, err)
				continue
			}
		}

		router := chi.NewRouter()
		router.Use(chain...)

		for _, route := range projectRoutes[project.ID] {
			upstream, ok := snapshot.Upstreams[project.ID][route.UpstreamName]
			if !ok {
				log.Printf("WARN: project %s route %s references unknown upstream %q, skipping",
					project.ID, route.Path, route.UpstreamName)
				continue
			}

			handler := NewProxyHandler(upstream, lbManager)
			var target chi.Router = router

			if route.AuthMode == models.AuthRequired {
				if jwtMW == nil {
					log.Printf("WARN: project %s route %s requires auth but no jwt middleware configured, skipping route",
						project.ID, route.Path)
					continue
				}
				target = router.With(jwtMW)
			}

			switch route.PathType {
			case models.PathExact:
				target.Method(route.Method, route.Path, handler)
			case models.PathPrefix:
				target.Method(route.Method, route.Path+"/*", handler)
			case models.PathRegex:
				target.Method(route.Method, route.Path, handler)
			}
		}

		snapshot.Projects[project.GatewayApiKey] = &RuntimeProject{
			ProjectID:   project.ID,
			Middlewares: project.Middlewares,
			Mux:         router,
		}
	}

	return nil
}

func NewProxyHandler(upstream *RuntimeUpstream, lbManager *loadbalancer.LBManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		backend, ok := loadbalancer.Pick(lbManager, upstream.Strategy, upstream.ID, upstream.Backends, r)
		if !ok {
			http.Error(w, "no healthy backend", http.StatusServiceUnavailable)
			return
		}

		lbManager.IncConn(upstream.ID, backend.Identifier())
		defer lbManager.DecConn(upstream.ID, backend.Identifier())

		backend.Proxy.ServeHTTP(w, r)
	}
}

func BuildMiddlewareChain(middlewares []models.Middleware, registry *pluginManager.PluginRegistry) ([]func(http.Handler) http.Handler, error) {

	sort.SliceStable(middlewares, func(i, j int) bool {
		oi, _ := middlewareOrder[middlewares[i].Name]
		oj, _ := middlewareOrder[middlewares[j].Name]
		return oi < oj
	})

	chain := make([]func(http.Handler) http.Handler, 0, len(middlewares))
	for _, mw := range middlewares {
		plugin, ok := registry.Get(mw.Name)
		if !ok {
			return nil, fmt.Errorf("middleware %q not registered", mw.Name)
		}
		mwFunc, err := plugin.CreateMiddleware(mw.Config)
		if err != nil {
			return nil, fmt.Errorf("middleware %q: %w", mw.Name, err)
		}
		chain = append(chain, mwFunc)
	}
	return chain, nil
}

func SplitJWT(middlewares []models.Middleware) ([]models.Middleware, *models.Middleware) {
	global := make([]models.Middleware, 0, len(middlewares))
	var jwt *models.Middleware
	for _, mw := range middlewares {
		if mw.Name == jwtPluginName {
			jwt = &mw
			continue
		}
		global = append(global, mw)
	}
	return global, jwt
}
