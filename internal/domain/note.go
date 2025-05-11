package domain

import (
	"time"
)

type Note struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Tags    []string  `json:"tags"`
	Created time.Time `json:"created"`
}
