package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Smil3MoreGH/gokeep/internal/handler"
	"github.com/Smil3MoreGH/gokeep/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	repo, err := storage.NewNoteRepository("notes.db")
	if err != nil {
		log.Fatalf("Fehler beim Öffnen der Datenbank: %v", err)
	}

	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/notes", handler.NewGetNotesHandler(repo))

	fmt.Println("Server läuft auf http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
