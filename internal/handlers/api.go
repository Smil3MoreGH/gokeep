// internal/handlers/api.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Smil3MoreGH/gokeep/internal/database"
	"github.com/Smil3MoreGH/gokeep/internal/models"
	"github.com/go-chi/chi/v5"
)

// APIHandler handles all API requests
type APIHandler struct {
	repo *database.NoteRepository
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(repo *database.NoteRepository) *APIHandler {
	return &APIHandler{repo: repo}
}

// GetAllNotes handles GET /api/notes
func (h *APIHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.repo.GetAll()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, notes)
}

// CreateNote handles POST /api/notes
func (h *APIHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.Create(&note); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, note)
}

// GetNote handles GET /api/notes/{id}
func (h *APIHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	note, err := h.repo.GetByID(id)
	if err != nil {
		if err.Error() == "note not found" {
			h.respondWithError(w, http.StatusNotFound, "Note not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, note)
}

// UpdateNote handles PUT /api/notes/{id}
func (h *APIHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	note.ID = id
	if err := h.repo.Update(&note); err != nil {
		if err.Error() == "note not found" {
			h.respondWithError(w, http.StatusNotFound, "Note not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, note)
}

// DeleteNote handles DELETE /api/notes/{id}
func (h *APIHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if err.Error() == "note not found" {
			h.respondWithError(w, http.StatusNotFound, "Note not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchNotes handles GET /api/notes/search?q=query
func (h *APIHandler) SearchNotes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	notes, err := h.repo.Search(query)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, notes)
}

// Helper methods

func (h *APIHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Error marshalling JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *APIHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}
