package handlers

import (
	"net/http"
	"quiz-app-be/internal/config"
	"quiz-app-be/internal/model"
	"quiz-app-be/internal/repository"
	"quiz-app-be/internal/service"
	"quiz-app-be/internal/setup/aws"

	"github.com/go-chi/chi"
	"github.com/go-chi/jsonp"
	"github.com/go-pg/pg/v10"
)

type Handler struct {
	userService   *service.UserService
	testService   *service.TestService
	resultService *service.ResultService
}

func NewHandler(
	pg *pg.DB,
	s3 *aws.AwsClient,
) *Handler {
	usersRepo := repository.NewUsers(pg)
	return &Handler{
		userService: service.NewUserService(usersRepo),
		testService: service.NewTestService(
			repository.NewTests(pg),
			usersRepo,
			s3,
		),
		resultService: service.NewResultService(
			repository.NewResults(pg),
			usersRepo,
		),
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

	router.Options("/create_test", h.createTest)
	router.Post("/create_test", h.createTest)

	router.Route("/test/{testID}", func(r chi.Router) {
		r.Options("/", h.getTest)
		r.Get("/", h.getTest)
	})

	router.Options("/set_result", h.setResult)
	router.Post("/set_result", h.setResult)

	router.Route("/get_home_tests/{number}", func(r chi.Router) {
		r.Options("/", h.getHomeTests)
		r.Get("/", h.getHomeTests)
	})

	router.Options("/count_all_public_tests", h.countAllPublicTests)
	router.Get("/count_all_public_tests", h.countAllPublicTests)

	router.Options("/my_profile", h.getProfile)
	router.Get("/my_profile", h.getProfile)

	router.Options("/my_results", h.getMyResults)
	router.Get("/my_results", h.getMyResults)

	router.Options("/my_tests_results", h.getMyTestsResults)
	router.Get("/my_tests_results", h.getMyTestsResults)

	router.Options("/my_tests", h.getMyTests)
	router.Get("/my_tests", h.getMyTests)

	router.Options("/get_all_tests", h.getAllTests)
	router.Get("/get_all_tests", h.getAllTests)

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
}
