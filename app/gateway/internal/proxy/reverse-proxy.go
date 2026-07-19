package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/lb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func BuildReverseProxy(upstreamID bson.ObjectID, target *url.URL, transport *http.Transport, lbManager *lb.LBManager) *httputil.ReverseProxy {
	backendKey := target.String()

	return &httputil.ReverseProxy{
		Transport: transport,
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(target)

			pr.Out.Header.Del("X-Gateway-Key")
		},

		ModifyResponse: func(resp *http.Response) error {
			if resp.StatusCode >= 500 {
				lbManager.MarkUnhealthy(upstreamID, backendKey)
			} else {
				lbManager.MarkHealthy(upstreamID, backendKey)
			}
			return nil
		},

		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			lbManager.MarkUnhealthy(upstreamID, backendKey)
			http.Error(w, "upstream unavailable", http.StatusBadGateway)
		},

		FlushInterval: -1,
	}
}
