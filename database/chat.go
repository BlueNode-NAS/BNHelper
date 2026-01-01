// Author: retrozenith <80767544+retrozenith@users.noreply.github.com>
// Description: Chat session and message storage

package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ChatSession struct {
	ID        int       `json:"id"`
	SessionID string    `json:"session_id"`
	Model     string    `json:"model"`
	Title     string    `json:"title,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChatMessage struct {
	ID        int       `json:"id"`
	SessionID string    `json:"session_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatStore struct {
	db *AIDB
}

func NewChatStore(db *AIDB) *ChatStore {
	return &ChatStore{db: db}
}

func (cs *ChatStore) CreateSession(model, title string) (*ChatSession, error) {
	sessionID := uuid.New().String()

	query := `
		INSERT INTO chat_sessions (session_id, model, title)
		VALUES (?, ?, ?)
	`

	result, err := cs.db.conn.Exec(query, sessionID, model, title)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat session: %w", err)
	}

	id, _ := result.LastInsertId()

	return &ChatSession{
		ID:        int(id),
		SessionID: sessionID,
		Model:     model,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (cs *ChatStore) GetSession(sessionID string) (*ChatSession, error) {
	query := `
		SELECT id, session_id, model, title, created_at, updated_at
		FROM chat_sessions
		WHERE session_id = ?
	`

	var session ChatSession
	err := cs.db.conn.QueryRow(query, sessionID).Scan(
		&session.ID,
		&session.SessionID,
		&session.Model,
		&session.Title,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("chat session not found: %s", sessionID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get chat session: %w", err)
	}

	return &session, nil
}

func (cs *ChatStore) ListSessions(limit int) ([]ChatSession, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `
		SELECT id, session_id, model, title, created_at, updated_at
		FROM chat_sessions
		ORDER BY updated_at DESC
		LIMIT ?
	`

	rows, err := cs.db.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list chat sessions: %w", err)
	}
	defer rows.Close()

	var sessions []ChatSession
	for rows.Next() {
		var session ChatSession
		err := rows.Scan(
			&session.ID,
			&session.SessionID,
			&session.Model,
			&session.Title,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chat session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

func (cs *ChatStore) DeleteSession(sessionID string) error {
	query := `DELETE FROM chat_sessions WHERE session_id = ?`

	result, err := cs.db.conn.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete chat session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat session not found: %s", sessionID)
	}

	return nil
}

func (cs *ChatStore) AddMessage(sessionID, role, content string) (*ChatMessage, error) {
	query := `
		INSERT INTO chat_messages (session_id, role, content)
		VALUES (?, ?, ?)
	`

	result, err := cs.db.conn.Exec(query, sessionID, role, content)
	if err != nil {
		return nil, fmt.Errorf("failed to add chat message: %w", err)
	}

	id, _ := result.LastInsertId()

	updateQuery := `UPDATE chat_sessions SET updated_at = CURRENT_TIMESTAMP WHERE session_id = ?`
	cs.db.conn.Exec(updateQuery, sessionID)

	return &ChatMessage{
		ID:        int(id),
		SessionID: sessionID,
		Role:      role,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}

func (cs *ChatStore) GetMessages(sessionID string) ([]ChatMessage, error) {
	query := `
		SELECT id, session_id, role, content, created_at
		FROM chat_messages
		WHERE session_id = ?
		ORDER BY created_at ASC
	`

	rows, err := cs.db.conn.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		err := rows.Scan(
			&msg.ID,
			&msg.SessionID,
			&msg.Role,
			&msg.Content,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chat message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (cs *ChatStore) UpdateSessionTitle(sessionID, title string) error {
	query := `UPDATE chat_sessions SET title = ? WHERE session_id = ?`

	result, err := cs.db.conn.Exec(query, title, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session title: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat session not found: %s", sessionID)
	}

	return nil
}

func float64SliceToBytes(slice []float64) ([]byte, error) {
	return json.Marshal(slice)
}

func bytesToFloat64Slice(data []byte) ([]float64, error) {
	var slice []float64
	err := json.Unmarshal(data, &slice)
	return slice, err
}
