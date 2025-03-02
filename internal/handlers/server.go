package handlers

import (
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	mux := http.NewServeMux()

	// Регистрируем обработчики
	mux.HandleFunc("/create_test", createTest)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/refresh", refresh)

	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	return &Server{
		httpServer: server,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
