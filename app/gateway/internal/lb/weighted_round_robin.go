package lb

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// AI genrated

func pickWeightedRoundRobin[T Backend](
	manager *LBManager,
	upstreamID bson.ObjectID,
	backends []T,
) (zero T, _ bool) {
	if len(backends) == 0 {
		return zero, false
	}

	state := manager.getOrCreate(upstreamID)

	// SWRR updates multiple backend states atomically.
	state.mu.Lock()
	defer state.mu.Unlock()

	var (
		best        T
		bestState   *backendState
		totalWeight int64
		found       bool
	)

	for _, backend := range backends {
		bs := state.state(backend.Identifier())

		if !bs.healthy.Load() {
			continue
		}

		weight := int64(backend.GetWeight())
		if weight <= 0 {
			weight = 1
		}

		totalWeight += weight

		current := bs.currentWeight.Add(weight)

		if !found || current > bestState.currentWeight.Load() {
			best = backend
			bestState = bs
			found = true
		}
	}

	if !found {
		return zero, false
	}

	bestState.currentWeight.Add(-totalWeight)

	return best, true
}
