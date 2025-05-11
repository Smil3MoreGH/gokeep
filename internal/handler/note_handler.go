package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Smil3MoreGH/gokeep/internal/storage"
)

func NewGetNotesHandler(repo *storage.NoteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notes, err := repo.FindAll()
		if err != nil {
			http.Error(w, "Fehler beim Abrufen der Notizen", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
	}
}
