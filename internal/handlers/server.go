package handlers

import (
	"net/http"
	"quiz-app-be/internal/config"
	"quiz-app-be/internal/model"
	"quiz-app-be/internal/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/jsonp"
)

type Handler struct {
	userService *service.UserService
}

func NewHandler(userService *service.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) Routing(cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()
	router.Use(jsonp.Handler)

	router.Options("/login", h.login)
	router.Post("/login", h.login)
	router.Options("/refresh", h.refresh)
	router.Post("/refresh", h.refresh)
	router.Options("/register", h.register)
	router.Post("/register", h.register)
	router.Post("/create_test", h.createTest)

	return router
}

func setError(w http.ResponseWriter, err error, msg string) {
	switch err.Error() {
	case model.ErrEmptyFields:
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	case model.ErrInvalidPassword:
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	case model.ErrUserNotFound:
		http.Error(w, "User not found", http.StatusNotFound)
		return
	case model.ErrUserAlreadyExists:
		http.Error(w, "User already exists", http.StatusConflict)
		return
	case model.ErrGenerateTokens:
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	case model.ErrInvalidRefreshToken:
		http.Error(w, "Refresh token not validated", http.StatusUnauthorized)
		return
	default:
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
	}
	return
}
