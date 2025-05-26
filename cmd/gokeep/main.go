// cmd/gokeep/main.go
package gokeep

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

//go:embed web/*
var webFS embed.FS

func main() {
	// Initialize the app
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// Setup routes
	setupRoutes(r)

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

func setupRoutes(r chi.Router) {
	// go-app handler
	r.Handle("/", &app.Handler{
		Name:        "Gokeep",
		Description: "A minimalist note-taking app built with Go",
		Styles: []string{
			"/web/app.css",
		},
		Scripts: []string{
			"/web/app.js",
		},
		CacheableResources: []string{
			"/web/app.css",
			"/web/app.js",
		},
	})

	// Static files
	r.Handle("/web/*", http.StripPrefix("/web/", http.FileServer(http.FS(webFS))))

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))

		r.Route("/notes", func(r chi.Router) {
			r.Get("/", getAllNotes)
			r.Post("/", createNote)
			r.Get("/search", searchNotes)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", getNote)
				r.Put("/", updateNote)
				r.Delete("/", deleteNote)
			})
		})
	})
}

// Placeholder handlers - to be implemented
func getAllNotes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"notes": []}`)
}

func createNote(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, `{"id": 1, "title": "New Note", "content": ""}`)
}

func getNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Fprintf(w, `{"id": %s, "title": "Note", "content": "Content"}`, id)
}

func updateNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Fprintf(w, `{"id": %s, "title": "Updated Note", "content": "Updated"}`, id)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func searchNotes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	fmt.Fprintf(w, `{"query": "%s", "notes": []}`, query)
}
