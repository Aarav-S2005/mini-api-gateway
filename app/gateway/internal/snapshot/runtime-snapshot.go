package snapshot

import (
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RuntimeSnapshot struct {
	Projects     map[string]*Project // apiKey -> Project
	ProjectMuxes map[string]http.Handler
	Upstreams    map[bson.ObjectID]map[string]*Upstream // ProjectID -> upstreamName -> Upstream
}
