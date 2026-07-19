package sync

import (
	"context"
	"net/http"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/lb"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/store"
	pluginManager "github.com/Aarav-S2005/mini-api-gateway/app/plugin-manager/registry"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SynchronizeLoader struct {
	registry       *store.Registry
	db             *mongo.Database
	pluginRegistry *pluginManager.PluginRegistry

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
	case config.ResourceProject:
		return sl.loadProject(ctx, notification)

	case config.ResourceRoute:
		return sl.loadRoute(ctx, notification)

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

func (sl *SynchronizeLoader) loadProject(ctx context.Context, n config.UpdateEventNotification) error {
	return sl.updateSnapshot(func(snap *store.Snapshot) error {
		return nil
	})
}
func (sl *SynchronizeLoader) loadUpstream(ctx context.Context, n config.UpdateEventNotification) error {
	return sl.updateSnapshot(func(snap *store.Snapshot) error {
		return nil
	})
}
func (sl *SynchronizeLoader) loadRoute(ctx context.Context, n config.UpdateEventNotification) error {
	return sl.updateSnapshot(func(snap *store.Snapshot) error {
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
