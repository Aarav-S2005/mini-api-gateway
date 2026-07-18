package store

import (
	"crypto/rsa"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Aarav-S2005/mini-api-gateway/app/db/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Snapshot struct {
	Projects  map[string]*RuntimeProject                    // map api-gateway-key to project
	Upstreams map[bson.ObjectID]map[string]*RuntimeUpstream // map projectID to upstreams
}

type RuntimeProject struct {
	ProjectID   bson.ObjectID
	JWTKey      *rsa.PublicKey
	Middlewares []models.Middleware
	Mux         http.Handler
}

type RuntimeUpstream struct {
	ID       bson.ObjectID
	Strategy models.LoadBalancingStrategy
	Backends []RuntimeBackend
}

type RuntimeBackend struct {
	URL    *url.URL
	Proxy  *httputil.ReverseProxy
	Weight int
}

func (b RuntimeBackend) Identifier() string { return b.URL.String() }
func (b RuntimeBackend) GetWeight() int     { return b.Weight }
