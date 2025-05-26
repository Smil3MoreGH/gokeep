// internal/ui/app.go
package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Smil3MoreGH/gokeep/internal/models"
	"github.com/Smil3MoreGH/gokeep/internal/ui/components"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type App struct {
	app.Compo

	notes         []models.Note
	searchTerm    string
	isLoading     bool
	error         error
	editingNoteID int64
	newNote       models.Note
	showNewNote   bool
}

func (a *App) OnMount(ctx app.Context) {
	a.loadNotes(ctx)
}

func (a *App) Render() app.UI {
	return app.Div().Class("app-container").Body(
		a.renderHeader(),
		app.Div().Class("main-content").Body(
			a.renderNewNote(),
			a.renderError(),
			app.If(
				a.isLoading,
				func() app.UI {
					return app.Div().Class("loading").Text("Loading notes...")
				},
			),
			app.If(
				!a.isLoading && len(a.notes) > 0,
				func() app.UI {
					return a.renderNotesGrid()
				},
			),
			app.If(
				!a.isLoading && len(a.notes) == 0 && a.searchTerm == "",
				func() app.UI {
					return app.Div().Class("empty-state").Body(
						app.H2().Text("No notes yet"),
						app.P().Text("Click the + button to create your first note"),
					)
				},
			),
		),
	)
}

func (a *App) renderHeader() app.UI {
	return app.Header().Class("app-header").Body(
		app.Div().Class("header-content").Body(
			app.H1().Class("app-title").Text("Gokeep"),
			app.Div().Class("search-container").Body(
				app.Input().
					Type("search").
					Class("search-input").
					Placeholder("Search notes...").
					Value(a.searchTerm).
					OnInput(a.onSearchInput),
			),
		),
	)
}

func (a *App) renderNewNote() app.UI {
	if !a.showNewNote {
		return app.Button().
			Class("fab").
			Title("Create new note").
			OnClick(a.onNewNoteClick).
			Body(app.Text("+"))
	}

	return app.Div().Class("new-note-container").Body(
		app.Div().Class("note-card new-note").Body(
			app.Input().
				Type("text").
				Class("note-title-input").
				Placeholder("Title").
				Value(a.newNote.Title).
				OnInput(a.onNewNoteTitleInput).
				AutoFocus(true),
			app.Textarea().
				Class("note-content-input").
				Placeholder("Take a note...").
				Rows(3).
				Text(a.newNote.Content).
				On("input", a.onNewNoteContentInput),
			app.Div().Class("note-actions").Body(
				app.Button().
					Class("btn btn-primary").
					Text("Save").
					OnClick(a.onSaveNewNote).
					Disabled(a.newNote.Title == "" && a.newNote.Content == ""),
				app.Button().
					Class("btn btn-secondary").
					Text("Cancel").
					OnClick(a.onCancelNewNote),
			),
		),
	)
}

func (a *App) renderNotesGrid() app.UI {
	filteredNotes := a.getFilteredNotes()
	return app.Div().Class("notes-grid").Body(
		app.Range(filteredNotes).Slice(func(i int) app.UI {
			note := filteredNotes[i]
			// v10: Direkt als *components.NoteCard – ist korrekt
			return &components.NoteCard{
				Note:      note,
				IsEditing: note.ID == a.editingNoteID,
				OnEdit:    a.onEditNote,
				OnDelete:  a.onDeleteNote,
				OnSave:    a.onSaveNote,
				OnCancel:  a.onCancelEdit,
			}
		}),
	)
}

// renderError renders error messages
func (a *App) renderError() app.UI {
	if a.error == nil {
		return nil
	}

	return app.Div().Class("error-message").Body(
		app.Text(a.error.Error()),
		app.Button().
			Class("close-error").
			Text("×").
			OnClick(a.onCloseError),
	)
}

// Event Handlers

func (a *App) onSearchInput(ctx app.Context, e app.Event) {
	a.searchTerm = ctx.JSSrc().Get("value").String()
	ctx.Update()
}

func (a *App) onNewNoteClick(ctx app.Context, e app.Event) {
	a.showNewNote = true
	a.newNote = models.Note{}
	ctx.Update()
}

func (a *App) onNewNoteTitleInput(ctx app.Context, e app.Event) {
	a.newNote.Title = ctx.JSSrc().Get("value").String()
	ctx.Update()
}

func (a *App) onNewNoteContentInput(ctx app.Context, e app.Event) {
	a.newNote.Content = ctx.JSSrc().Get("value").String()
	ctx.Update()
}

func (a *App) onSaveNewNote(ctx app.Context, e app.Event) {
	a.createNote(ctx)
}

func (a *App) onCancelNewNote(ctx app.Context, e app.Event) {
	a.showNewNote = false
	a.newNote = models.Note{}
	ctx.Update()
}

func (a *App) onEditNote(ctx app.Context, noteID int64) {
	a.editingNoteID = noteID
	ctx.Update()
}

func (a *App) onDeleteNote(ctx app.Context, noteID int64) {
	a.deleteNote(ctx, noteID)
}

func (a *App) onSaveNote(ctx app.Context, note models.Note) {
	a.updateNote(ctx, note)
}

func (a *App) onCancelEdit(ctx app.Context) {
	a.editingNoteID = 0
	ctx.Update()
}

func (a *App) onCloseError(ctx app.Context, e app.Event) {
	a.error = nil
	ctx.Update()
}

// Helper Methods

func (a *App) getFilteredNotes() []models.Note {
	if a.searchTerm == "" {
		return a.notes
	}

	filtered := make([]models.Note, 0)
	for _, note := range a.notes {
		if contains(note.Title, a.searchTerm) || contains(note.Content, a.searchTerm) {
			filtered = append(filtered, note)
		}
	}
	return filtered
}

// API Methods

func (a *App) loadNotes(ctx app.Context) {
	a.isLoading = true
	ctx.Update()

	go func() {
		resp, err := http.Get("/api/notes")
		if err != nil {
			a.error = err
			a.isLoading = false
			ctx.Update()
			return
		}
		defer resp.Body.Close()

		var notes []models.Note
		if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
			a.error = err
			a.isLoading = false
			ctx.Update()
			return
		}

		a.notes = notes
		a.isLoading = false
		ctx.Update()
	}()
}

func (a *App) createNote(ctx app.Context) {
	noteJSON, err := json.Marshal(a.newNote)
	if err != nil {
		a.error = err
		ctx.Update()
		return
	}

	go func() {
		resp, err := http.Post("/api/notes", "application/json", bytes.NewReader(noteJSON))
		if err != nil {
			a.error = err
			ctx.Update()
			return
		}
		defer resp.Body.Close()

		var createdNote models.Note
		if err := json.NewDecoder(resp.Body).Decode(&createdNote); err != nil {
			a.error = err
			ctx.Update()
			return
		}

		a.notes = append([]models.Note{createdNote}, a.notes...)
		a.showNewNote = false
		a.newNote = models.Note{}
		ctx.Update()
	}()
}

func (a *App) updateNote(ctx app.Context, note models.Note) {
	noteJSON, err := json.Marshal(note)
	if err != nil {
		a.error = err
		ctx.Update()
		return
	}

	go func() {
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/notes/%d", note.ID), bytes.NewReader(noteJSON))
		if err != nil {
			a.error = err
			a.editingNoteID = 0
			ctx.Dispatch(func(ctx app.Context) {
				ctx.Update()
			})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			a.error = err
			a.editingNoteID = 0
			ctx.Dispatch(func(ctx app.Context) {
				ctx.Update()
			})
			return
		}
		defer resp.Body.Close()

		// Update local state
		for i, n := range a.notes {
			if n.ID == note.ID {
				a.notes[i] = note
				break
			}
		}
		a.editingNoteID = 0
		ctx.Dispatch(func(ctx app.Context) {
			ctx.Update()
		})
	}()
}

func (a *App) deleteNote(ctx app.Context, noteID int64) {
	go func() {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/notes/%d", noteID), nil)
		if err != nil {
			a.error = err
			ctx.Dispatch(func(ctx app.Context) {
				ctx.Update()
			})
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			a.error = err
			ctx.Dispatch(func(ctx app.Context) {
				ctx.Update()
			})
			return
		}
		defer resp.Body.Close()

		// Remove from local state
		filtered := make([]models.Note, 0)
		for _, note := range a.notes {
			if note.ID != noteID {
				filtered = append(filtered, note)
			}
		}
		a.notes = filtered
		ctx.Dispatch(func(ctx app.Context) {
			ctx.Update()
		})
	}()
}

// Utility function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
