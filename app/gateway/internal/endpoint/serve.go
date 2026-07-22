package endpoint

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/store"
)

const keyHeader = "X-Gateway-Key"
const serverPrefix = "/api"

type Handler struct {
	snapshotRegistry *store.Registry
}

func NewHandler(snapshotRegistry *store.Registry) *Handler {
	return &Handler{snapshotRegistry: snapshotRegistry}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	apiKey := r.Header.Get(keyHeader)
	if apiKey == "" {
		http.Error(w, "missing gateway key header: "+keyHeader, http.StatusUnauthorized)
		return
	}

	snap := h.snapshotRegistry
	project, ok := snap.Get().Projects[apiKey]
	if !ok {
		http.Error(w, "invalid gateway key", http.StatusUnauthorized)
		return
	}
	r.URL.Path = strings.TrimPrefix(r.URL.Path, serverPrefix)
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}
	project.Mux.ServeHTTP(w, r)
	log.Printf("%s %s project=%s status_took=%s", r.Method, r.URL.Path, project.ProjectID.Hex(), time.Since(start))
}
