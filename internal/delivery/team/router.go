package team

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/team/add", h.Add)
	r.Get("/team/get", h.Get)
}
