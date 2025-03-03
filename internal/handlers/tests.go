package handlers

import (
	"encoding/json"
	"net/http"
	modelhttp "quiz-app-be/internal/model/modelHttp"

	"github.com/go-chi/chi"
)

func (h *Handler) createTest(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req modelhttp.CreateTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.testService.CreateTest(req)
	if err != nil {
		setError(w, err, "Failed to create test: ")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getTest(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	testID := chi.URLParam(r, "testID")
	resp, err := h.testService.GetTest(testID)
	if err != nil {
		setError(w, err, "Failed to get test: ")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
