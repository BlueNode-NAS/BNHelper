// Author: retrozenith <80767544+retrozenith@users.noreply.github.com>
// Description: Database connection and initialization

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
	defaultDBPath = "/var/lib/bnhelper/bnhelper.db"
	defaultDBDir  = "/var/lib/bnhelper"
)

type DB struct {
	conn *sql.DB
}

func New(dbPath string) (*DB, error) {
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	conn.SetMaxOpenConns(1)

	db := &DB{conn: conn}

	if err := db.initialize(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	log.Printf("Database initialized at: %s", dbPath)
	return db, nil
}

func (db *DB) initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS configurations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT NOT NULL UNIQUE,
		value TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_configurations_key ON configurations(key);

	CREATE TRIGGER IF NOT EXISTS update_configurations_timestamp 
	AFTER UPDATE ON configurations
	FOR EACH ROW
	BEGIN
		UPDATE configurations SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END;
	`

	_, err := db.conn.Exec(schema)
	return err
}

func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func (db *DB) GetConnection() *sql.DB {
	return db.conn
}
