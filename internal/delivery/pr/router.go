package pr

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/pullRequest/create", h.Create)
	r.Post("/pullRequest/merge", h.Merge)
	r.Post("/pullRequest/reassign", h.Reassign)
}
