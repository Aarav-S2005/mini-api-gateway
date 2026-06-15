package project

import (
	"net/http"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/app_error"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/lib"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/middleware"
	"github.com/Aarav-S2005/mini-api-gateway/app/plugin-manager/registry"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
	reg     *registry.PluginRegistry
	pub     *config.Publisher
}

func NewHandler(repo *Repository, pluginRegistry registry.PluginRegistry, pub *config.Publisher) chi.Router {
	h := &Handler{service: NewService(repo), reg: &pluginRegistry, pub: pub}
	return h.initRoutes()
}

func (h *Handler) initRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(
		middleware.CookieToBearer,
		middleware.Verifier(),
		middleware.Authenticator(),
	)
	r.Group(func(r chi.Router) {
		r.Post("/", h.createProject) // done
		r.Get("/", h.getAllProjects) // done
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.getProject)                              // done
			r.Delete("/", h.deleteProject)                        // done
			r.Patch("/middleware", h.updateMiddlewares)           // done
			r.Patch("/accesslist", h.updateAccessList)            // done
			r.Patch("/loadbalancer", h.updateLoadBalancerConfig)  // done
			r.Delete("/middleware/{name}", h.deleteMiddleware)    // done
			r.Delete("/loadbalancer", h.deleteLoadBalancerConfig) // done

			// Route Handlers
			r.Route("/routes", func(r chi.Router) {
				r.Get("/", h.getProjectRoutes) // done
				r.Post("/add", h.addRoute)
				r.Route("/{routeID}", func(r chi.Router) {
					r.Patch("/update", h.updateRoute)
					r.Delete("/delete", h.DeleteRoute)
				})
			})
		})
	})
	return r
}

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserID(r.Context())
	var req CreatProjectRequest
	err = lib.ConvertJSONToStruct(r, &req)
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("could not parse json", err))
		return
	}
	res, err := h.service.createProject(r.Context(), req, userID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	lib.ConvertStructToJSON(w, http.StatusCreated, res)
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	projectID, userID, err := GetProjectAndUserID(r)
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("could not get userID or projectID", err))
		return
	}
	res, err := h.service.GetProject(r.Context(), projectID, userID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	lib.ConvertStructToJSON(w, http.StatusOK, res)
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	projectID, userID, err := GetProjectAndUserID(r)
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("could not get userID or projectID", err))
		return
	}
	err = h.service.DeleteProject(r.Context(), projectID, userID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getAllProjects(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserID(r.Context())
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("", err))
		return
	}
	res, err := h.service.GetAllProjects(r.Context(), userID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	lib.ConvertStructToJSON(w, http.StatusOK, res)
}

func (h *Handler) updateMiddlewares(w http.ResponseWriter, r *http.Request) {
	projectID, userID, err := GetProjectAndUserID(r)
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("", err))
		return
	}
	var req UpdateMiddlewaresRequest
	err = h.service.UpdateMiddlewares(r.Context(), projectID, userID, req, *h.reg)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateAccessList(w http.ResponseWriter, r *http.Request) {
	projectID, userID, err := GetProjectAndUserID(r)
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("no userID", err))
		return
	}
	var req UpdateAccessListRequest
	err = lib.ConvertJSONToStruct(r, &req)
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("could not parse json", err))
		return
	}
	err = h.service.UpdateAccessList(r.Context(), projectID, userID, req)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateLoadBalancerConfig(w http.ResponseWriter, r *http.Request) {
	var req UpdateLoadBalancerConfigRequest
	err := lib.ConvertJSONToStruct(r, &req)
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("could not parse json", err))
		return
	}

}

func (h *Handler) deleteMiddleware(w http.ResponseWriter, r *http.Request) {
	mwname := chi.URLParam(r, "name")
	projectID, userID, err := GetProjectAndUserID(r)
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("", err))
		return
	}
	err = h.service.deleteMiddleware(r.Context(), projectID, userID, mwname)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteLoadBalancerConfig(w http.ResponseWriter, r *http.Request) {
	projectID, userID, err := GetProjectAndUserID(r)
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("", err))
		return
	}
	err = h.service.deleteLoadBalancerConfig(r.Context(), projectID, userID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ROUTE HANDLER FUNCTIONS

func (h *Handler) getProjectRoutes(w http.ResponseWriter, r *http.Request) {
	projectID, userID, err := GetProjectAndUserID(r)
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("", err))
		return
	}
	routes, err := h.service.getProjectRoutes(r.Context(), projectID, userID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	lib.ConvertStructToJSON(w, http.StatusOK, routes)
}

func (h *Handler) addRoute(w http.ResponseWriter, r *http.Request) {
	var addRouteRequest AddUpdateRouteRequest
	err := lib.ConvertJSONToStruct(r, &addRouteRequest)
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("could not parse json", err))
		return
	}
	projectID, userID, err := GetProjectAndUserID(r)
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("", err))
		return
	}
	id, err := h.service.addProjectRoute(r.Context(), projectID, userID, addRouteRequest)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	lib.ConvertStructToJSON(w, http.StatusOK, AddRouteResponse{ID: id})
}

func (h *Handler) updateRoute(w http.ResponseWriter, r *http.Request) {
	var req AddUpdateRouteRequest
	err := lib.ConvertJSONToStruct(r, &req)
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("could not parse json", err))
		return
	}
	projectID, userID, err := GetProjectAndUserID(r)
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("", err))
		return
	}
	routeID, err := GetIdFromEndpoint(r, "routeID")
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("no routeID", err))
		return
	}
	err = h.service.updateProjectRoute(r.Context(), routeID, projectID, userID, req)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteRoute(w http.ResponseWriter, r *http.Request) {
	projectID, userID, err := GetProjectAndUserID(r)
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("", err))
		return
	}
	routeID, err := GetIdFromEndpoint(r, "routeID")
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("no routeID", err))
		return
	}
	err = h.service.deleteProjectRoute(r.Context(), routeID, projectID, userID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
