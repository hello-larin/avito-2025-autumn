package pr

import (
	"encoding/jsogithub.com/hello-larin/avito-2025-autumn
	"net/http"

	customerror "avito-2025-autumn/internal/error"

	"github.com/go-playground/validator/v10"
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req prCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON_PARSE_ERROR", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, "JSON_PARSE_ERROR", http.StatusBadRequest)
		return
	}

	pr := req.ToModel()

	createdPR, reviewers, err := h.useCase.CreatePullRequest(r.Context(), &pr)
	if err != nil {
		customerror.WriteErrorResponse(w, err)
		return
	}

	resp := toCreateResponse(createdPR, reviewers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "JSON_ENCODE_ERROR", http.StatusBadRequest)
	}
}

func (h *Handler) Merge(w http.ResponseWriter, r *http.Request) {
	var req prMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON_PARSE_ERROR", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, "JSON_PARSE_ERROR", http.StatusBadRequest)
		return
	}

	mergedPR, reviewers, err := h.useCase.MergePullRequest(r.Context(), req.PullRequestID)
	if err != nil {
		customerror.WriteErrorResponse(w, err)
		return
	}

	resp := toMergeResponse(mergedPR, reviewers)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "JSON_ENCODE_ERROR", http.StatusBadRequest)
	}
}

func (h *Handler) Reassign(w http.ResponseWriter, r *http.Request) {
	var req prReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON_PARSE_ERROR", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, "JSON_PARSE_ERROR", http.StatusBadRequest)
		return
	}

	pr, reviewers, newReviewerID, err := h.useCase.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		customerror.WriteErrorResponse(w, err)
		return
	}

	resp := toReassignResponse(pr, reviewers, newReviewerID)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "JSON_ENCODE_ERROR", http.StatusBadRequest)
	}
}
