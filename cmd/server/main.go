package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Smil3MoreGH/gokeep/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Health Check
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Neue Route: /notes
	r.Get("/notes", handler.GetNotesHandler)

	fmt.Println("Server läuft auf http://localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
