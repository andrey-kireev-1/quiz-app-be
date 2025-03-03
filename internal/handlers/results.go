package handlers

import (
	"encoding/json"
	"net/http"
	modelhttp "quiz-app-be/internal/model/modelHttp"
)

func (h *Handler) setResult(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req modelhttp.SetResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.resultService.SetResult(req)
	if err != nil {
		setError(w, err, "Failed to set result: ")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getMyResults(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		http.Error(w, "No refresh token provided", http.StatusUnauthorized)
		return
	}
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	results, err := h.resultService.GetUserResults(accessToken)
	if err != nil {
		setError(w, err, "Failed to get results: ")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) getMyTestsResults(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		http.Error(w, "No refresh token provided", http.StatusUnauthorized)
		return
	}
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	results, err := h.resultService.GetAuthorTestsResults(accessToken)
	if err != nil {
		setError(w, err, "Failed to get results: ")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
