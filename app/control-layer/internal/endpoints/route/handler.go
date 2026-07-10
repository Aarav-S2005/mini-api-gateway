package route

import (
	"net/http"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/app_error"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/lib"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
	pub     *config.Publisher
}

func NewHandler(repo *Repository, pub *config.Publisher) *Handler {
	return &Handler{service: NewService(repo), pub: pub}
}

// REMEMBER: Redis Pub left when config changes

func (h *Handler) initRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.createRoute)
	r.Get("/{routeID}", h.getRoute)
	r.Get("/", h.getAllRoutes)
	r.Put("/{routeID}", h.updateRoute)
	r.Delete("/{routeID}", h.deleteRoute)

	return r
}

func (h *Handler) createRoute(w http.ResponseWriter, r *http.Request) {
	var reqBody CreateOrUpdateRouteRequestDTO
	err := lib.ConvertJSONToStruct(r, &reqBody)
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("could not parse json", err))
		return
	}
	userID, projectID, err := getUserIDAndProjectID(r)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	res, err := h.service.createRoute(r.Context(), reqBody, userID, projectID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	lib.ConvertStructToJSON(w, 201, res)
}

func (h *Handler) getRoute(w http.ResponseWriter, r *http.Request) {
	userID, projectID, routeID, err := getUserIDProjectIDAndRouteID(r)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	route, err := h.service.getRouteByID(r.Context(), userID, projectID, routeID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	lib.ConvertStructToJSON(w, 200, route)
}

func (h *Handler) getAllRoutes(w http.ResponseWriter, r *http.Request) {
	userID, projectID, err := getUserIDAndProjectID(r)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	routes, err := h.service.getAllRoute(r.Context(), userID, projectID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	lib.ConvertStructToJSON(w, 200, routes)
}

func (h *Handler) updateRoute(w http.ResponseWriter, r *http.Request) {
	userID, projectID, routeID, err := getUserIDProjectIDAndRouteID(r)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	var reqBody CreateOrUpdateRouteRequestDTO
	err = lib.ConvertJSONToStruct(r, &reqBody)
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("could not parse json", err))
		return
	}
	err = h.service.updateRoute(r.Context(), reqBody, userID, projectID, routeID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteRoute(w http.ResponseWriter, r *http.Request) {
	userID, projectID, routeID, err := getUserIDProjectIDAndRouteID(r)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	err = h.service.deleteRoute(r.Context(), userID, projectID, routeID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
