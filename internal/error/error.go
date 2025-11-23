package customerror

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

var (
	ErrTeamExists  = errors.New("TEAM_EXISTS")
	ErrPRExists    = errors.New("PR_EXISTS")
	ErrPRMerged    = errors.New("PR_MERGED")
	ErrNotAssigned = errors.New("NOT_ASSIGNED")
	ErrNoCandidate = errors.New("NO_CANDIDATE")
	ErrNotFound    = errors.New("NOT_FOUND")
	ErrJSONParse   = errors.New("JSON_ERROR")
)

func WriteErrorResponse(w http.ResponseWriter, err error) {
	var errorMap = map[string]struct {
		statusCode int
		message    string
	}{
		ErrTeamExists.Error():  {http.StatusConflict, "team_name already exists"},
		ErrPRExists.Error():    {http.StatusConflict, "PR id already exists"},
		ErrPRMerged.Error():    {http.StatusConflict, "cannot reassign on merged PR"},
		ErrNotAssigned.Error(): {http.StatusConflict, "reviewer is not assigned to this PR"},
		ErrNoCandidate.Error(): {http.StatusConflict, "no active replacement candidate in team"},
		ErrNotFound.Error():    {http.StatusNotFound, "resource not found"},
		ErrJSONParse.Error():   {http.StatusBadRequest, "cannot parse json"},
	}
	info, exists := errorMap[err.Error()]
	if !exists {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    err.Error(),
			Message: info.message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(info.statusCode)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON_ENCODE_ERROR", http.StatusBadRequest)
	}
}
