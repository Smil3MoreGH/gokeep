// internal/database/repository.go
package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Smil3MoreGH/gokeep/internal/models"
)

// NoteRepository handles all database operations for notes
type NoteRepository struct {
	db *DB
}

// NewNoteRepository creates a new note repository
func NewNoteRepository(db *DB) *NoteRepository {
	return &NoteRepository{db: db}
}

// Create inserts a new note into the database
func (r *NoteRepository) Create(note *models.Note) error {
	note.SetDefaults()

	query := `
        INSERT INTO notes (title, content, color, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
        RETURNING id
    `

	err := r.db.conn.QueryRow(
		query,
		note.Title,
		note.Content,
		note.Color,
		note.CreatedAt,
		note.UpdatedAt,
	).Scan(&note.ID)

	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	return nil
}

// GetAll retrieves all notes from the database
func (r *NoteRepository) GetAll() ([]models.Note, error) {
	query := `
        SELECT id, title, content, color, created_at, updated_at
        FROM notes
        ORDER BY updated_at DESC
    `

	rows, err := r.db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.Color,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	return notes, rows.Err()
}

// GetByID retrieves a single note by its ID
func (r *NoteRepository) GetByID(id int64) (*models.Note, error) {
	query := `
        SELECT id, title, content, color, created_at, updated_at
        FROM notes
        WHERE id = ?
    `

	var note models.Note
	err := r.db.conn.QueryRow(query, id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.Color,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("note not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	return &note, nil
}

// Update updates an existing note
func (r *NoteRepository) Update(note *models.Note) error {
	note.UpdatedAt = time.Now()

	query := `
        UPDATE notes 
        SET title = ?, content = ?, color = ?, updated_at = ?
        WHERE id = ?
    `

	result, err := r.db.conn.Exec(
		query,
		note.Title,
		note.Content,
		note.Color,
		note.UpdatedAt,
		note.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}

// Delete removes a note from the database
func (r *NoteRepository) Delete(id int64) error {
	query := `DELETE FROM notes WHERE id = ?`

	result, err := r.db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}

// Search performs a full-text search on notes
func (r *NoteRepository) Search(query string) ([]models.Note, error) {
	// Clean and prepare search query
	searchQuery := strings.TrimSpace(query)
	if searchQuery == "" {
		return r.GetAll()
	}

	// Use FTS5 for search
	sqlQuery := `
        SELECT n.id, n.title, n.content, n.color, n.created_at, n.updated_at
        FROM notes n
        JOIN notes_fts fts ON n.id = fts.rowid
        WHERE notes_fts MATCH ?
        ORDER BY rank
    `

	rows, err := r.db.conn.Query(sqlQuery, searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to search notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.Color,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	return notes, rows.Err()
}

// Count returns the total number of notes
func (r *NoteRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notes`

	err := r.db.conn.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count notes: %w", err)
	}

	return count, nil
}
