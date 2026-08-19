package auth_transport_http

import "github.com/go-chi/chi/v5"

func (h *AuthHTTPHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.LoginUser)
	})
}
