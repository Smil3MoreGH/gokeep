package handler

import (
	"encoding/json"
	"github.com/Smil3MoreGH/gokeep/internal/domain"
	"github.com/google/uuid"
	"net/http"
	"time"

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

func NewPostNoteHandler(repo *storage.NoteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var note domain.Note

		// JSON Body einlesen
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			http.Error(w, "Ungültige Eingabe", http.StatusBadRequest)
			return
		}

		// ID + Timestamp generieren
		note.ID = uuid.New().String()
		note.Created = time.Now()

		// In DB speichern
		if err := repo.Save(note); err != nil {
			http.Error(w, "Fehler beim Speichern", http.StatusInternalServerError)
			return
		}

		// Erfolgsantwort
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}
