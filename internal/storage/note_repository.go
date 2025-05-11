package storage

import (
	"encoding/json"
	"errors"

	"github.com/Smil3MoreGH/gokeep/internal/domain"
	"go.etcd.io/bbolt"
)

var notesBucket = []byte("notes")

type NoteRepository struct {
	db *bbolt.DB
}

func NewNoteRepository(path string) (*NoteRepository, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	// Bucket anlegen, wenn nicht vorhanden
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(notesBucket)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &NoteRepository{db: db}, nil
}

func (r *NoteRepository) Save(note domain.Note) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(notesBucket)
		encoded, err := json.Marshal(note)
		if err != nil {
			return err
		}
		return b.Put([]byte(note.ID), encoded)
	})
}

func (r *NoteRepository) FindAll() ([]domain.Note, error) {
	var notes []domain.Note

	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(notesBucket)
		if b == nil {
			return errors.New("notes bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var note domain.Note
			if err := json.Unmarshal(v, &note); err != nil {
				return err
			}
			notes = append(notes, note)
			return nil
		})
	})

	return notes, err
}
