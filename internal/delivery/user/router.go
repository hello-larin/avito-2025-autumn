package user

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/users/setIsActive", h.SetActive)
	r.Get("/users/getReview", h.GetReview)
}
