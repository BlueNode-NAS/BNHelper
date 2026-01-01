// Author: retrozenith <80767544+retrozenith@users.noreply.github.com>
// Description: Configuration storage and retrieval

package database

import (
	"database/sql"
	"fmt"
	"time"
)

type Configuration struct {
	ID          int       `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ConfigStore struct {
	db *DB
}

func NewConfigStore(db *DB) *ConfigStore {
	return &ConfigStore{db: db}
}

func (cs *ConfigStore) Set(key, value, description string) error {
	query := `
		INSERT INTO configurations (key, value, description)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			description = excluded.description,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := cs.db.conn.Exec(query, key, value, description)
	if err != nil {
		return fmt.Errorf("failed to set configuration: %w", err)
	}
	return nil
}

func (cs *ConfigStore) Get(key string) (*Configuration, error) {
	query := `
		SELECT id, key, value, description, created_at, updated_at
		FROM configurations
		WHERE key = ?
	`

	var config Configuration
	err := cs.db.conn.QueryRow(query, key).Scan(
		&config.ID,
		&config.Key,
		&config.Value,
		&config.Description,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("configuration not found: %s", key)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	return &config, nil
}

func (cs *ConfigStore) GetAll() ([]Configuration, error) {
	query := `
		SELECT id, key, value, description, created_at, updated_at
		FROM configurations
		ORDER BY key
	`

	rows, err := cs.db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query configurations: %w", err)
	}
	defer rows.Close()

	var configs []Configuration
	for rows.Next() {
		var config Configuration
		err := rows.Scan(
			&config.ID,
			&config.Key,
			&config.Value,
			&config.Description,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan configuration: %w", err)
		}
		configs = append(configs, config)
	}

	return configs, rows.Err()
}

func (cs *ConfigStore) Delete(key string) error {
	query := `DELETE FROM configurations WHERE key = ?`

	result, err := cs.db.conn.Exec(query, key)
	if err != nil {
		return fmt.Errorf("failed to delete configuration: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("configuration not found: %s", key)
	}

	return nil
}

func (cs *ConfigStore) Exists(key string) (bool, error) {
	query := `SELECT COUNT(*) FROM configurations WHERE key = ?`

	var count int
	err := cs.db.conn.QueryRow(query, key).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check configuration existence: %w", err)
	}

	return count > 0, nil
}
