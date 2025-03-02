package handlers

import (
	"net/http"
)

func (h *Handler) createTest(w http.ResponseWriter, r *http.Request) {
	allowCors(w)
	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}
