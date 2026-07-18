package lb

import (
	"math/rand/v2"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func pickRandom[T Backend](
	manager *LBManager,
	upstreamID bson.ObjectID,
	backends []T,
) (zero T, _ bool) {
	n := len(backends)
	if n == 0 {
		return zero, false
	}

	state := manager.getOrCreate(upstreamID)

	start := rand.IntN(n)

	for i := 0; i < n; i++ {
		backend := backends[(start+i)%n]

		if state.state(backend.Identifier()).healthy.Load() {
			return backend, true
		}
	}

	return zero, false
}
