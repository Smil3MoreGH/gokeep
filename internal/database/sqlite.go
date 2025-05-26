// internal/database/sqlite.go
package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB represents the database connection
type DB struct {
	conn *sql.DB
}

// NewDB creates a new database connection
func NewDB(dataSourceName string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}

	// Run migrations
	if err := db.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Migrate creates the necessary tables
func (db *DB) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS notes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        content TEXT,
        color TEXT DEFAULT '#ffffff',
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at);
    CREATE INDEX IF NOT EXISTS idx_notes_updated_at ON notes(updated_at);
    
    -- Full-text search table
    CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
        title, 
        content, 
        content_rowid=id
    );

    -- Triggers to keep FTS table in sync
    CREATE TRIGGER IF NOT EXISTS notes_ai AFTER INSERT ON notes
    BEGIN
        INSERT INTO notes_fts(rowid, title, content) 
        VALUES (new.id, new.title, new.content);
    END;

    CREATE TRIGGER IF NOT EXISTS notes_ad AFTER DELETE ON notes
    BEGIN
        DELETE FROM notes_fts WHERE rowid = old.id;
    END;

    CREATE TRIGGER IF NOT EXISTS notes_au AFTER UPDATE ON notes
    BEGIN
        UPDATE notes_fts 
        SET title = new.title, content = new.content 
        WHERE rowid = new.id;
    END;
    `

	_, err := db.conn.Exec(query)
	return err
}

// BeginTx starts a new transaction
func (db *DB) BeginTx() (*sql.Tx, error) {
	return db.conn.Begin()
}
