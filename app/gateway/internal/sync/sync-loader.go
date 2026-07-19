package sync

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config"
	"github.com/Aarav-S2005/mini-api-gateway/app/db/models"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/lb"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/store"
	pluginManager "github.com/Aarav-S2005/mini-api-gateway/app/plugin-manager/registry"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SynchronizeLoader struct {
	registry *store.Registry
	db       *mongo.Database

	transport *http.Transport
	lbManager *lb.LBManager
	plugins   *pluginManager.PluginRegistry
}

func NewSynchronizeLoader(registry *store.Registry, db *mongo.Database, transport *http.Transport, lbManager *lb.LBManager, plugins *pluginManager.PluginRegistry) *SynchronizeLoader {
	return &SynchronizeLoader{
		registry:  registry,
		db:        db,
		transport: transport,
		lbManager: lbManager,
		plugins:   plugins,
	}
}

func (sl *SynchronizeLoader) handler(ctx context.Context, notification config.UpdateEventNotification) error {
	switch notification.Resource {
	case config.ResourceProject, config.ResourceRoute:
		return sl.loadRuntimeProject(ctx, notification)
	case config.ResourceUpstream:
		return sl.loadUpstream(ctx, notification)
	default:
		return nil
	}
}

func (sl *SynchronizeLoader) updateSnapshot(update func(*store.Snapshot) error) error {
	old := sl.registry.Get()
	snap := cloneSnapshot(old)
	if err := update(snap); err != nil {
		return err
	}
	sl.registry.Store(snap)
	return nil
}

func (sl *SynchronizeLoader) loadRuntimeProject(ctx context.Context, n config.UpdateEventNotification) error {
	return sl.updateSnapshot(func(snap *store.Snapshot) error {
		var project models.Project
		var projectID bson.ObjectID
		if n.Resource == config.ResourceRoute {
			var route models.Route
			err := sl.db.Collection("routes").FindOne(ctx, bson.M{"_id": n.ResourceID}).Decode(&route)
			if err != nil {
				// i have not checked if route is deleted because if route is deleted,
				// redis will send that project has been updated and hence resource type is project
				return fmt.Errorf("could not load Route: %s due to %w", n.ResourceID.String(), err)
			}
			projectID = route.ProjectID
		} else {
			projectID = n.ResourceID
		}
		err := sl.db.Collection("projects").FindOne(ctx, bson.M{"_id": projectID}).Decode(&project)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil
			}
			return fmt.Errorf("could not load RuntimeProject: %s due to %w", projectID.String(), err)
		}
		var routes []models.Route
		cursor, err := sl.db.Collection("routes").Find(ctx, bson.M{"project_id": projectID, "enabled": true})
		if err != nil {
			return fmt.Errorf("could not get routes from project %s due to %w", projectID.String(), err)
		}
		defer cursor.Close(ctx)
		if err := cursor.All(ctx, &routes); err != nil {
			return fmt.Errorf("could not get routes from project %s due to %w", projectID.String(), err)
		}
		globalMiddlewares, jwtConfig := store.SplitJWT(project.Middlewares)
		chain, err := store.BuildMiddlewareChain(globalMiddlewares, sl.plugins)
		if err != nil {
			return fmt.Errorf("ERROR: project %s has invalid middleware config, excluding from snapshot: %v", project.ID, err)
		}

		var jwtMW func(http.Handler) http.Handler
		if jwtConfig != nil {
			plugin, ok := sl.plugins.Get(jwtConfig.Name)
			if !ok {
				return fmt.Errorf("project %s references unknown jwt plugin %q", project.ID, jwtConfig.Name)
			}
			jwtMW, err = plugin.CreateMiddleware(jwtConfig.Config)
			if err != nil {
				return fmt.Errorf("ERROR: project %s jwt middleware config invalid, excluding from snapshot: %v", project.ID, err)
			}
		}
		router := chi.NewRouter()
		router.Use(chain...)
		for _, route := range routes {
			upstream, ok := snap.Upstreams[project.ID][route.UpstreamName]
			if !ok {
				log.Printf("WARN: project %s route %s references unknown upstream %q, skipping", project.ID, route.Path, route.UpstreamName)
				continue
			}
			handler := store.NewProxyHandler(upstream, sl.lbManager)
			var target chi.Router = router
			if route.AuthMode == models.AuthRequired {
				if jwtMW == nil {
					log.Printf("WARN: project %s route %s requires auth but no jwt middleware configured, skipping route", project.ID, route.Path)
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
		snap.Projects[project.GatewayApiKey] = &store.RuntimeProject{
			ProjectID:   project.ID,
			Middlewares: project.Middlewares,
			Mux:         router,
		}
		return nil
	})
}
func (sl *SynchronizeLoader) loadUpstream(ctx context.Context, n config.UpdateEventNotification) error {
	return sl.updateSnapshot(func(snap *store.Snapshot) error {
		var upstream models.Upstream
		err := sl.db.Collection("upstreams").FindOne(ctx, bson.M{"_id": n.ResourceID}).Decode(&upstream)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				delete(snap.Upstreams, n.ResourceID)
				return nil
			}
			return fmt.Errorf("could not load Upstream: %s due to %w", n.ResourceID.String(), err)
		}
		err = store.LoadSingleRuntimeUpstream(upstream, snap, sl.transport, sl.lbManager)
		if err != nil {
			return fmt.Errorf("could not load Upstream: %s due to %w", n.ResourceID.String(), err)
		}
		return nil
	})
}

func cloneSnapshot(old *store.Snapshot) *store.Snapshot {
	snap := &store.Snapshot{
		Projects:  make(map[string]*store.RuntimeProject, len(old.Projects)),
		Upstreams: make(map[bson.ObjectID]map[string]*store.RuntimeUpstream, len(old.Upstreams)),
	}

	for k, v := range old.Projects {
		snap.Projects[k] = v
	}

	for pid, ups := range old.Upstreams {
		m := make(map[string]*store.RuntimeUpstream, len(ups))
		for name, rt := range ups {
			m[name] = rt
		}
		snap.Upstreams[pid] = m
	}

	return snap
}
