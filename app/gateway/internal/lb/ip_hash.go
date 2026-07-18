package lb

import (
	"hash/fnv"
	"net"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func pickIPHash[T Backend](
	manager *LBManager,
	upstreamID bson.ObjectID,
	backends []T,
	r *http.Request,
) (zero T, _ bool) {
	n := len(backends)
	if n == 0 {
		return zero, false
	}

	state := manager.getOrCreate(upstreamID)

	start := hash(clientIP(r)) % uint64(n)

	for i := 0; i < n; i++ {
		backend := backends[(start+uint64(i))%uint64(n)]

		if state.state(backend.Identifier()).healthy.Load() {
			return backend, true
		}
	}

	return zero, false
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if idx := strings.IndexByte(xff, ','); idx >= 0 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}

func hash(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}
