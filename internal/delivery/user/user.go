package user

import (
	"encoding/json"
	"net/http"

	customerror "github.com/hello-larin/avito-2025-autumn/internal/error"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	useCase  usecase
	validate *validator.Validate
}

func New(uc usecase, validate *validator.Validate) *Handler {
	return &Handler{
		useCase:  uc,
		validate: validator.New(),
	}
}

func (h *Handler) SetActive(w http.ResponseWriter, r *http.Request) {
	var req userSetActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON_PARSE_ERROR", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, "JSON_PARSE_ERROR", http.StatusBadRequest)
		return
	}

	user, err := h.useCase.SetUserActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		customerror.WriteErrorResponse(w, err)
		return
	}

	resp := toSetActiveResponse(user)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "JSON_ENCODE_ERROR", http.StatusBadRequest)
	}
}

func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "MISSING_PARAM", http.StatusBadRequest)
		return
	}

	result, err := h.useCase.GetUserAssignedPRs(r.Context(), userID)
	if err != nil {
		customerror.WriteErrorResponse(w, err)
		return
	}

	resp := toPRResponse(userID, result)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "JSON_ENCODE_ERROR", http.StatusBadRequest)
	}
}
