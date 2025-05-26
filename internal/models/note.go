// internal/models/note.go
package models

import (
	"time"
)

// Note represents a single note in the application
type Note struct {
	ID        int64     `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Color     string    `json:"color" db:"color"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NoteColor represents available note colors
type NoteColor string

const (
	ColorWhite  NoteColor = "#ffffff"
	ColorYellow NoteColor = "#fff475"
	ColorOrange NoteColor = "#fbbc04"
	ColorPink   NoteColor = "#f28b82"
	ColorPurple NoteColor = "#d7aefb"
	ColorBlue   NoteColor = "#aecbfa"
	ColorGreen  NoteColor = "#ccff90"
	ColorGray   NoteColor = "#e8eaed"
)

// ValidateColor checks if the provided color is valid
func ValidateColor(color string) bool {
	validColors := []NoteColor{
		ColorWhite, ColorYellow, ColorOrange, ColorPink,
		ColorPurple, ColorBlue, ColorGreen, ColorGray,
	}

	for _, c := range validColors {
		if string(c) == color {
			return true
		}
	}
	return false
}

// SetDefaults sets default values for a new note
func (n *Note) SetDefaults() {
	if n.Color == "" {
		n.Color = string(ColorWhite)
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	if n.UpdatedAt.IsZero() {
		n.UpdatedAt = time.Now()
	}
}

// Update updates the note with new values
func (n *Note) Update(title, content, color string) {
	if title != "" {
		n.Title = title
	}
	if content != "" {
		n.Content = content
	}
	if color != "" && ValidateColor(color) {
		n.Color = color
	}
	n.UpdatedAt = time.Now()
}
