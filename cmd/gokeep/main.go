package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/maxence-charriere/go-app/v10/pkg/app"

	"github.com/Smil3MoreGH/gokeep/internal/database"
	"github.com/Smil3MoreGH/gokeep/internal/handlers"
	"github.com/Smil3MoreGH/gokeep/internal/ui"
)

//go:embed web/*
var webFS embed.FS

func main() {
	// Initialise SQLite database (creates file if it does not exist)
	files, err := webFS.ReadDir("web")
	if err != nil {
		log.Fatal("webFS.ReadDir:", err)
	}
	for _, f := range files {
		log.Println("Embedded file:", f.Name())
	}

	db, err := database.NewDB("gokeep.db")
	if err != nil {
		log.Fatalf("failed to initialise database: %v", err)
	}
	defer db.Close()

	// Repository & REST handler layer
	repo := database.NewNoteRepository(db)
	api := handlers.NewAPIHandler(repo)

	// Register UI route for client‑side Go‑app components when running in the browser
	app.Route("/", func() app.Composer { return &ui.App{} })
	app.RunWhenOnBrowser()

	// Router / middleware stack
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Serve the UI (root path)
	r.Handle("/", &app.Handler{
		Name:        "Gokeep",
		Title:       "Gokeep",
		Description: "A minimalist note‑taking app written 100 % in Go",
		Styles: []string{
			"/web/app.css", // optional – remove if you do not ship CSS yet
		},
		Scripts: []string{
			"/web/app.js", // optional – remove if you do not ship JS yet
		},
		CacheableResources: []string{
			"/web/app.css",
			"/web/app.js",
		},
	})

	// Serve static assets that the Go‑app bundle references (favicon, CSS, etc.)
	r.Handle("/web/*", http.FileServer(http.FS(webFS)))

	// Wire up JSON API underneath /api
	setupAPIRoutes(r, api)

	// HTTP server with graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Run server in background goroutine so we can listen for OS signals
	go func() {
		log.Printf("Gokeep listening on http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server stopped unexpectedly: %v", err)
		}
	}()

	// Wait for Ctrl‑C / SIGTERM
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Println("shutdown signal received – stopping …")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

// setupAPIRoutes registers /api/... endpoints backed by the API handler.
func setupAPIRoutes(r chi.Router, h *handlers.APIHandler) {
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content‑Type", "application/json"))

		r.Route("/notes", func(r chi.Router) {
			r.Get("/", h.GetAllNotes)
			r.Post("/", h.CreateNote)

			// Search endpoint: /api/notes/search?q=foo
			r.Get("/search", h.SearchNotes)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetNote)
				r.Put("/", h.UpdateNote)
				r.Delete("/", h.DeleteNote)
			})
		})
	})
}
