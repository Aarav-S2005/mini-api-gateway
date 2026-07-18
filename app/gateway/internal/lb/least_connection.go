package lb

import (
	"math"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func pickLeastConnections[T Backend](
	manager *LBManager,
	upstreamID bson.ObjectID,
	backends []T,
) (zero T, _ bool) {
	if len(backends) == 0 {
		return zero, false
	}

	state := manager.getOrCreate(upstreamID)

	var (
		best      T
		bestConns int64 = math.MaxInt64
		found     bool
	)

	for _, backend := range backends {
		backend_State := state.state(backend.Identifier())

		if !backend_State.healthy.Load() {
			continue
		}

		conns := backend_State.activeConn.Load()

		if !found || conns < bestConns {
			best = backend
			bestConns = conns
			found = true
		}
	}

	return best, found
}
