package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Smil3MoreGH/gokeep/internal/domain"
	"github.com/google/uuid"
)

func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	// Beispielnotizen (später ersetzen wir das durch DB-Zugriff)
	notes := []domain.Note{
		{
			ID:      uuid.New().String(),
			Title:   "Erste Notiz",
			Content: "Das ist eine Beispielnotiz.",
			Tags:    []string{"example", "test"},
			Created: time.Now(),
		},
		{
			ID:      uuid.New().String(),
			Title:   "Zweite Notiz",
			Content: "Noch eine Notiz, um die Liste zu füllen.",
			Tags:    []string{"todo"},
			Created: time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}
