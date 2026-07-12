package snapshot

import "github.com/Aarav-S2005/mini-api-gateway/app/db/models"

type Route struct {
	Path         string
	PathType     models.PathType
	Method       string
	UpstreamName string
	AuthMode     models.AuthMode
	Enabled      bool
}
