// Author: retrozenith <80767544+retrozenith@users.noreply.github.com>
// Description: AI database connection and initialization

package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	defaultAIDBPath = "/var/lib/bnhelper/ai.db"
)

type AIDB struct {
	conn *sql.DB
}

func NewAIDB(dbPath string) (*AIDB, error) {
	if dbPath == "" {
		dbPath = defaultAIDBPath
	}

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create AI database directory: %w", err)
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open AI database: %w", err)
	}

	conn.SetMaxOpenConns(1)

	db := &AIDB{conn: conn}

	if err := db.initialize(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize AI database: %w", err)
	}

	log.Printf("AI database initialized at: %s", dbPath)
	return db, nil
}

func (db *AIDB) initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS chat_sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL UNIQUE,
		model TEXT NOT NULL,
		title TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS chat_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (session_id) REFERENCES chat_sessions(session_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS indexed_files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT NOT NULL UNIQUE,
		content TEXT NOT NULL,
		embedding BLOB,
		model TEXT NOT NULL,
		file_size INTEGER,
		file_hash TEXT,
		indexed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS file_chunks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_id INTEGER NOT NULL,
		chunk_index INTEGER NOT NULL,
		content TEXT NOT NULL,
		embedding BLOB,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (file_id) REFERENCES indexed_files(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_chat_sessions_session_id ON chat_sessions(session_id);
	CREATE INDEX IF NOT EXISTS idx_chat_messages_session_id ON chat_messages(session_id);
	CREATE INDEX IF NOT EXISTS idx_indexed_files_path ON indexed_files(file_path);
	CREATE INDEX IF NOT EXISTS idx_indexed_files_hash ON indexed_files(file_hash);
	CREATE INDEX IF NOT EXISTS idx_file_chunks_file_id ON file_chunks(file_id);

	CREATE TRIGGER IF NOT EXISTS update_chat_sessions_timestamp 
	AFTER UPDATE ON chat_sessions
	FOR EACH ROW
	BEGIN
		UPDATE chat_sessions SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END;

	CREATE TRIGGER IF NOT EXISTS update_indexed_files_timestamp 
	AFTER UPDATE ON indexed_files
	FOR EACH ROW
	BEGIN
		UPDATE indexed_files SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END;
	`

	_, err := db.conn.Exec(schema)
	return err
}

func (db *AIDB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func (db *AIDB) GetConnection() *sql.DB {
	return db.conn
}
