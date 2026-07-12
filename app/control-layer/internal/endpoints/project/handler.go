package project

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/app_error"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/lib"
	"github.com/Aarav-S2005/mini-api-gateway/app/plugin-manager/registry"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
	reg     *registry.PluginRegistry
}

func NewHandler(repo *Repository, pluginRegistry *registry.PluginRegistry, upstreamHandler, routeHandler chi.Router) chi.Router {
	h := &Handler{service: NewService(repo), reg: pluginRegistry}
	return h.initRoutes(upstreamHandler, routeHandler)
}

func (h *Handler) initRoutes(upstreamHandler, routeHandler chi.Router) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/", h.createProject) // done
		r.Get("/", h.getAllProjects) // done
		r.Route("/{projectID}", func(r chi.Router) {
			r.Get("/", h.getProject)                            // done
			r.Delete("/", h.deleteProject)                      // done
			r.Patch("/middlewares", h.updateMiddlewares)        // done
			r.Patch("/accesslist", h.updateAccessList)          // done
			r.Delete("/middlewares/{name}", h.deleteMiddleware) // done
			r.Mount("/upstreams", upstreamHandler)
			r.Mount("/routes", routeHandler)
		})
	})
	return r
}

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserID(r.Context())
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("failed to get userID", err))
		return
	}
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
	projectID, userID, err := lib.GetProjectAndUserID(r)
	fmt.Println(projectID, userID)
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
	projectID, userID, err := lib.GetProjectAndUserID(r)
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
		app_error.HandleError(w, app_error.Unauthorized("failed to get userID", err))
		return
	}
	res, err := h.service.GetAllProjects(r.Context(), userID)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	lib.ConvertStructToJSON(w, http.StatusOK, GetAllProjectResponse{Projects: res})
}

func (h *Handler) updateMiddlewares(w http.ResponseWriter, r *http.Request) {
	projectID, userID, err := lib.GetProjectAndUserID(r)
	if err != nil {
		app_error.HandleError(w, app_error.Unauthorized("", err))
		return
	}
	var req UpdateMiddlewaresRequest
	err = lib.ConvertJSONToStruct(r, &req)
	if err != nil {
		app_error.HandleError(w, app_error.BadRequest("could not parse json", err))
		return
	}
	err = h.service.UpdateMiddlewares(r.Context(), projectID, userID, req, *h.reg)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateAccessList(w http.ResponseWriter, r *http.Request) {
	projectID, userID, err := lib.GetProjectAndUserID(r)
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

func (h *Handler) deleteMiddleware(w http.ResponseWriter, r *http.Request) {
	mwname := chi.URLParam(r, "name")
	if mwname == "" {
		app_error.HandleError(w, app_error.BadRequest("no middleware name", errors.New("middleware name is missing")))
		return
	}
	projectID, userID, err := lib.GetProjectAndUserID(r)
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
