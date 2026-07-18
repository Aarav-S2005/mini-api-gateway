package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/Aarav-S2005/mini-api-gateway/app/db/models"
	loadbalancer "github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/lb"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrCouldNotLoadFromDB = errors.New("could not load from database")

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

func LoadSnapshot(db *mongo.Database, transport *http.Transport, lbManager *loadbalancer.LBManager) (*Snapshot, error) {
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
	if err = loadProjects(projects, routes, snapshot, lbManager); err != nil {
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
				Proxy:  buildReverseProxy(upstream.ID, target, transport, lbManager),
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

func loadProjects(projects []models.Project, routes []models.Route, snapshot *Snapshot, lbManager *loadbalancer.LBManager) error {
	projectRoutes := make(map[bson.ObjectID][]models.Route)
	for _, route := range routes {
		if !route.Enabled {
			continue
		}

		projectRoutes[route.ProjectID] = append(projectRoutes[route.ProjectID], route)
	}

	for _, project := range projects {

		router := chi.NewRouter()

		// will build middleware chain

		for _, route := range projectRoutes[project.ID] {
			upstream, ok := snapshot.Upstreams[project.ID][route.UpstreamName]
			if !ok {
				log.Printf("WARN: project %s route %s references unknown upstream %q, skipping",
					project.ID, route.Path, route.UpstreamName)
				continue
			}

			handler := NewProxyHandler(upstream, lbManager)

			switch route.PathType {
			case models.PathExact:
				router.Method(route.Method, route.Path, handler)
			case models.PathPrefix:
				router.Method(route.Method, route.Path+"/*", handler)
			case models.PathRegex:
				router.Method(route.Method, route.Path, handler)
			}
		}

		snapshot.Projects[project.GatewayApiKey] = &RuntimeProject{
			ProjectID:   project.ID,
			JWTKey:      nil, // will load if needed only, else it will be included in middlewrae chain
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

func buildReverseProxy(upstreamID bson.ObjectID, target *url.URL, transport *http.Transport, lbManager *loadbalancer.LBManager) *httputil.ReverseProxy {
	backendKey := target.String()

	return &httputil.ReverseProxy{
		Transport: transport,
		Director: func(r *http.Request) {
			r.URL.Scheme = target.Scheme
			r.URL.Host = target.Host
			r.Host = target.Host
			r.Header.Del("X-Gateway-Key")
		},

		// ModifyResponse sees actual HTTP responses from the backend —
		// this is where 5xx detection belongs, not ErrorHandler.
		ModifyResponse: func(resp *http.Response) error {
			if resp.StatusCode >= 500 {
				lbManager.MarkUnhealthy(upstreamID, backendKey)
			} else {
				// reactive recovery: a good response means the backend is alive.
				// Active health checks (once wired) will do this more reliably;
				// this just avoids a backend staying marked-down forever between checks.
				lbManager.MarkHealthy(upstreamID, backendKey)
			}
			return nil
		},

		// ErrorHandler only fires on transport-level failures — dial refused,
		// timeout, TLS handshake failure, connection reset. No HTTP response
		// exists yet at this point, so this is always "unreachable", not "5xx".
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			lbManager.MarkUnhealthy(upstreamID, backendKey)
			http.Error(w, "upstream unavailable", http.StatusBadGateway)
		},

		FlushInterval: -1,
	}
}
