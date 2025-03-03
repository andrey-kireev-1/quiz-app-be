package handlers

import (
	"encoding/json"
	"net/http"
	modelhttp "quiz-app-be/internal/model/modelHttp"
)

func allowCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, credentials")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq modelhttp.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := h.userService.LoginUser(loginReq)
	if err != nil {
		setError(w, err, "Failed to login user: ")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var refreshReq modelhttp.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := h.userService.RefreshToken(refreshReq)
	if err != nil {
		setError(w, err, "Failed to refresh token: ")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var registerReq modelhttp.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := h.userService.Register(registerReq)
	if err != nil {
		setError(w, err, "Failed to register user: ")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tokens)
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get refresh token from header
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		http.Error(w, "No refresh token provided", http.StatusUnauthorized)
		return
	}
	// Remove "Bearer " prefix if present
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	// Get user profile from service
	profile, err := h.userService.GetProfile(accessToken)
	if err != nil {
		setError(w, err, "Failed to get profile: ")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
