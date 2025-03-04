package handlers

import (
	"encoding/json"
	"net/http"
	modelhttp "quiz-app-be/internal/model/modelHttp"
	"time"

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

func (h *Handler) getHomeTests(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	paginationNum := chi.URLParam(r, "number")
	resp, err := h.testService.GetHomeTests(paginationNum)
	if err != nil {
		setError(w, err, "Failed to get test: ")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) countAllPublicTests(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resp, err := h.testService.CountAllPublicTests()
	if err != nil {
		setError(w, err, "Failed to get test: ")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) getMyTests(w http.ResponseWriter, r *http.Request) {
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

	tests, err := h.testService.GetUserTests(accessToken)
	if err != nil {
		setError(w, err, "Failed to get tests: ")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tests)
}

func (h *Handler) getAllTests(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	filters := modelhttp.TestFilters{}

	if testName := r.URL.Query().Get("testName"); testName != "" {
		filters.TestName = &testName
	}

	if authorName := r.URL.Query().Get("authorName"); authorName != "" {
		filters.AuthorName = &authorName
	}

	if createdFrom := r.URL.Query().Get("createdFrom"); createdFrom != "" {
		t, err := time.Parse(time.RFC3339, createdFrom)
		if err == nil {
			filters.CreatedFrom = &t
		}
	}

	if createdTo := r.URL.Query().Get("createdTo"); createdTo != "" {
		t, err := time.Parse(time.RFC3339, createdTo)
		if err == nil {
			filters.CreatedTo = &t
		}
	}

	if isStrict := r.URL.Query().Get("isStrict"); isStrict != "" {
		strict := isStrict == "true"
		filters.IsStrict = &strict
	}

	tests, err := h.testService.GetFilteredTests(filters)
	if err != nil {
		setError(w, err, "Failed to get tests: ")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tests)
}
