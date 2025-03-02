package main

import (
	"log"
	"quiz-app-be/internal/handlers"
)

func main() {
	srv := handlers.NewServer()
	log.Printf("Server starting on http://localhost:8000")
	log.Fatal(srv.Start())
}
