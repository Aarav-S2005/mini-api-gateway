package health

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/lb"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	checkInterval      = 10 * time.Second
	checkTimeout       = 3 * time.Second
	healthyThreshold   = 3
	unhealthyThreshold = 3
	defaultHealthPath  = "/health"
)

type counters struct {
	consecFail atomic.Int32
	consecOK   atomic.Int32
}

type Checker struct {
	httpClient *http.Client
	mu         sync.Mutex
	state      map[string]*counters
	registry   *store.Registry
	lbManager  *lb.LBManager
}

func NewChecker(transport *http.Transport, registry *store.Registry, lbManager *lb.LBManager) *Checker {
	return &Checker{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   checkTimeout,
		},
		state:     make(map[string]*counters),
		registry:  registry,
		lbManager: lbManager,
	}
}

func (c *Checker) Run(ctx context.Context) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.RunOnce(ctx)
		}
	}
}

func (c *Checker) RunOnce(ctx context.Context) {
	snap := c.registry.Get()
	var wg sync.WaitGroup
	for _, upstreams := range snap.Upstreams {
		for _, upstream := range upstreams {
			for _, backend := range upstream.Backends {
				wg.Add(1)
				go func(upstreamID bson.ObjectID, b store.RuntimeBackend) {
					defer wg.Done()
					c.check(ctx, upstreamID, b)
				}(upstream.ID, backend)
			}
		}
	}
	wg.Wait()
}

func (c *Checker) check(ctx context.Context, upstreamID bson.ObjectID, backend store.RuntimeBackend) {
	key := backend.Identifier()
	ok := c.pingHTTP(ctx, backend.URL) || c.pingTCP(backend.URL)
	counter := c.counterFor(key)
	if ok {
		counter.consecFail.Store(0)
		if counter.consecOK.Add(1) >= healthyThreshold {
			c.lbManager.MarkHealthy(upstreamID, key)
		}
	} else {
		counter.consecOK.Store(0)
		if counter.consecFail.Add(1) >= unhealthyThreshold {
			c.lbManager.MarkUnhealthy(upstreamID, key)
		}
	}
}

// Pingers
func (c *Checker) pingHTTP(ctx context.Context, target *url.URL) bool {
	healthURL := *target
	healthURL.Path = defaultHealthPath

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthURL.String(), nil)
	if err != nil {
		return false
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 500
}

func (c *Checker) pingTCP(target *url.URL) bool {
	conn, err := net.DialTimeout("tcp", target.Host, checkTimeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// HLEPER
func (c *Checker) counterFor(key string) *counters {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cnt, ok := c.state[key]; ok {
		return cnt
	}
	cnt := &counters{}
	c.state[key] = cnt
	return cnt
}
