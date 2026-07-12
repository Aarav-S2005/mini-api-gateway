package snapshot

import (
	"github.com/Aarav-S2005/mini-api-gateway/app/db/models"
)

type Upstream struct {
	LoadBalancingStrategy models.LoadBalancingStrategy
	Backends              []models.Backend
}
