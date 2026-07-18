package lb

import "go.mongodb.org/mongo-driver/v2/bson"

func pickRoundRobin[T Backend](
	manager *LBManager,
	upstreamID bson.ObjectID,
	backends []T,
) (zero T, _ bool) {
	n := len(backends)
	if n == 0 {
		return zero, false
	}

	state := manager.getOrCreate(upstreamID)

	start := (state.rrCounter.Add(1) - 1) % uint64(n)

	for i := 0; i < n; i++ {
		backend := backends[(start+uint64(i))%uint64(n)]

		if state.state(backend.Identifier()).healthy.Load() {
			return backend, true
		}
	}

	return zero, false
}
