package route

import (
	"net/http"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/app_error"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/lib"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
	pub     *config.Publisher
}

func NewHandler(repo *Repository, pub *config.Publisher) *Handler {
	return &Handler{service: NewService(repo), pub: pub}
}

func (h *Handler) initRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(
		middleware.CookieToBearer,
		middleware.Authenticator(),
		middleware.Verifier(),
	)
	r.Group(func(r chi.Router) {
		r.Get("/", h.getAllUpstream)
		r.Get("/{upstreamID}", h.getUpstream)
		r.Post("/", h.createUpstream)
		r.Put("/{upstreamID}", h.updateUpstream)
		r.Delete("/{upstreamID}", h.deleteUpstream)
	})
	return r
}

// redis pub left out
func (h *Handler) createUpstream(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserID(r.Context())
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("invalid token", err))
		return
	}
	projectID, err := lib.GetIdFromEndpoint(r, "projectID")
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("projectID missing from url", err))
		return
	}
	var reqBody CreateOrUpdateUpstreamRequestDTO
	err = lib.ConvertJSONToStruct(r, &reqBody)
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("invalid request body", err))
		return
	}
	upstreamID, err := h.service.createUpstream(r.Context(), userID, projectID, reqBody)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	lib.ConvertStructToJSON(w, 201, CreateUpstreamResponseDTO{UpstreamID: upstreamID})
}

func (h *Handler) getAllUpstream(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserID(r.Context())
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("invalid token", err))
		return
	}
	projectID, err := lib.GetIdFromEndpoint(r, "projectID")
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("projectID missing from url", err))
		return
	}
	upstreams, err := h.service.getAllUpstream(r.Context(), userID, projectID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	lib.ConvertStructToJSON(w, 200, upstreams)
}

func (h *Handler) getUpstream(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserID(r.Context())
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("invalid token", err))
		return
	}
	projectID, err := lib.GetIdFromEndpoint(r, "projectID")
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("projectID missing from url", err))
		return
	}
	upstreamID, err := lib.GetIdFromEndpoint(r, "upstreamID")
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("projectID missing from url", err))
		return
	}
	upstream, err := h.service.getUpstreamByID(r.Context(), userID, projectID, upstreamID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	lib.ConvertStructToJSON(w, 200, upstream)
}

// redis pub left out
func (h *Handler) updateUpstream(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserID(r.Context())
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("invalid token", err))
		return
	}
	projectID, err := lib.GetIdFromEndpoint(r, "projectID")
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("projectID missing from url", err))
		return
	}
	upstreamID, err := lib.GetIdFromEndpoint(r, "upstreamID")
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("projectID missing from url", err))
		return
	}
	var reqBody CreateOrUpdateUpstreamRequestDTO
	err = h.service.updateUpstream(r.Context(), userID, projectID, upstreamID, reqBody)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// redis pub left out
func (h *Handler) deleteUpstream(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserID(r.Context())
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("invalid token", err))
		return
	}
	projectID, err := lib.GetIdFromEndpoint(r, "projectID")
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("projectID missing from url", err))
		return
	}
	upstreamID, err := lib.GetIdFromEndpoint(r, "upstreamID")
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("projectID missing from url", err))
		return
	}
	err = h.service.deleteUpstream(r.Context(), userID, projectID, upstreamID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// finish later
func (h *Handler) publishChanges() {

}
