package lb

import (
	"sync"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type backendState struct {
	healthy       atomic.Bool
	activeConn    atomic.Int64
	currentWeight atomic.Int64
	failureCount  atomic.Int64
}

func newBackendState() *backendState {
	bs := &backendState{}
	bs.healthy.Store(true)
	return bs
}

type LBState struct {
	rrCounter atomic.Uint64
	backends  sync.Map // key: backend identifier -> *backendState
	mu        sync.Mutex
}

func (s *LBState) state(key string) *backendState {
	if v, ok := s.backends.Load(key); ok {
		return v.(*backendState)
	}
	actual, _ := s.backends.LoadOrStore(key, newBackendState())
	return actual.(*backendState)
}

type LBManager struct {
	mu     sync.RWMutex
	states map[bson.ObjectID]*LBState
}

func NewLBManager() *LBManager {
	return &LBManager{states: make(map[bson.ObjectID]*LBState)}
}

func (m *LBManager) getOrCreate(upstreamID bson.ObjectID) *LBState {
	m.mu.RLock()
	s, ok := m.states[upstreamID]
	m.mu.RUnlock()
	if ok {
		return s
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.states[upstreamID]; ok { // re-check under write lock
		return s
	}
	s = &LBState{}
	m.states[upstreamID] = s
	return s
}

func (m *LBManager) MarkUnhealthy(upstreamID bson.ObjectID, backendKey string) {
	state := m.getOrCreate(upstreamID).state(backendKey)

	if state.failureCount.Add(1) >= 10 {
		state.healthy.Store(false)
	}
}

func (m *LBManager) MarkHealthy(upstreamID bson.ObjectID, backendKey string) {
	state := m.getOrCreate(upstreamID).state(backendKey)

	state.failureCount.Store(0)
	state.healthy.Store(true)
}

func (m *LBManager) IncConn(upstreamID bson.ObjectID, backendKey string) {
	m.getOrCreate(upstreamID).state(backendKey).activeConn.Add(1)
}

func (m *LBManager) DecConn(upstreamID bson.ObjectID, backendKey string) {
	m.getOrCreate(upstreamID).state(backendKey).activeConn.Add(-1)
}

func (m *LBManager) GC(activeUpstreamIDs map[bson.ObjectID]struct{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id := range m.states {
		if _, ok := activeUpstreamIDs[id]; !ok {
			delete(m.states, id)
		}
	}
}
