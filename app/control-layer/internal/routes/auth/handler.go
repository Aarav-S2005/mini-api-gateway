package auth

import (
	"net/http"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/app_error"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/lib"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(repo *Repository) chi.Router {
	h := &Handler{service: NewService(repo)}
	r := h.initAuthRoutes()
	return r
}

func (h *Handler) initAuthRoutes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Post("/login", h.login)
		r.Post("/register", h.register)
		r.Group(func(r chi.Router) {
			r.Use(middleware.CookieToBearer)
			r.Use(middleware.Verifier())
			r.Use(middleware.Authenticator())
			r.Get("/logout", h.logout)
		})
	})
	return r
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestDTO
	err := lib.ConvertJSONToStruct(r, &reqBody)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	token, err := h.service.loginUser(r.Context(), reqBody)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	SetCookie(w, "auth-token", token)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	RemoveCookie(w, "auth-token")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestDTO
	err := lib.ConvertJSONToStruct(r, &reqBody)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	token, err := h.service.signUpUser(r.Context(), reqBody)
	if err != nil {
		app_error.HandleError(w, err)
		return
	}
	SetCookie(w, "x-auth-token", token)
	w.WriteHeader(http.StatusNoContent)
}
