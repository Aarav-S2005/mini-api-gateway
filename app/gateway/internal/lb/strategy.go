package lb

import (
	"net/http"

	"github.com/Aarav-S2005/mini-api-gateway/app/db/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Pick[T Backend](
	manager *LBManager,
	strategy models.LoadBalancingStrategy,
	upstreamID bson.ObjectID,
	backends []T,
	r *http.Request,
) (zero T, _ bool) {

	switch strategy {

	case models.RoundRobinLoadBalancing:
		return pickRoundRobin(manager, upstreamID, backends)

	case models.WeightedRoundRobin:
		return pickWeightedRoundRobin(manager, upstreamID, backends)

	case models.RandomLoadBalancing:
		return pickRandom(manager, upstreamID, backends)

	case models.IPHashLoadBalancing:
		return pickIPHash(manager, upstreamID, backends, r)

	case models.LeastConnections:
		return pickLeastConnections(manager, upstreamID, backends)

	default:
		return pickRoundRobin(manager, upstreamID, backends)
	}
}
