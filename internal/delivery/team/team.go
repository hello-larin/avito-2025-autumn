package team

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	customerror "github.com/hello-larin/avito-2025-autumn/internal/error"
	"github.com/hello-larin/avito-2025-autumn/internal/models"
)

type Handler struct {
	useCase  usecase
	validate *validator.Validate
}

func New(uc usecase, validate *validator.Validate) *Handler {
	return &Handler{
		useCase:  uc,
		validate: validate,
	}
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	var req team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON_PARSE_ERROR", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, "JSON_PARSE_ERROR", http.StatusBadRequest)
		return
	}

	var members []models.UserDB
	for _, m := range req.Members {
		members = append(members, m.toModel())
	}

	team, teamMembers, err := h.useCase.CreateTeam(r.Context(), req.TeamName, members)
	if err != nil {
		customerror.WriteErrorResponse(w, err)
		return
	}

	resp := toCreateResponse(team, teamMembers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "JSON_ENCODE_ERROR", http.StatusBadRequest)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		http.Error(w, "MISSING_PARAM", http.StatusBadRequest)
		return
	}

	team, members, err := h.useCase.GetTeam(r.Context(), teamName)
	if err != nil {
		customerror.WriteErrorResponse(w, err)
		return
	}

	resp := toDTO(team, members)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "JSON_ENCODE_ERROR", http.StatusBadRequest)
	}
}
