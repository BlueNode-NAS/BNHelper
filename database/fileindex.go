// Author: retrozenith <80767544+retrozenith@users.noreply.github.com>
// Description: File indexing storage and retrieval

package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"
)

type IndexedFile struct {
	ID        int       `json:"id"`
	FilePath  string    `json:"file_path"`
	Content   string    `json:"content"`
	Embedding []float64 `json:"embedding,omitempty"`
	Model     string    `json:"model"`
	FileSize  int64     `json:"file_size"`
	FileHash  string    `json:"file_hash"`
	IndexedAt time.Time `json:"indexed_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileChunk struct {
	ID         int       `json:"id"`
	FileID     int       `json:"file_id"`
	ChunkIndex int       `json:"chunk_index"`
	Content    string    `json:"content"`
	Embedding  []float64 `json:"embedding,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type FileIndexStore struct {
	db *AIDB
}

func NewFileIndexStore(db *AIDB) *FileIndexStore {
	return &FileIndexStore{db: db}
}

func (fis *FileIndexStore) IndexFile(filePath, content, model string, embedding []float64) (*IndexedFile, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	fileHash, err := calculateFileHash(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate file hash: %w", err)
	}

	embeddingBytes, err := float64SliceToBytes(embedding)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize embedding: %w", err)
	}

	query := `
		INSERT INTO indexed_files (file_path, content, embedding, model, file_size, file_hash)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(file_path) DO UPDATE SET
			content = excluded.content,
			embedding = excluded.embedding,
			model = excluded.model,
			file_size = excluded.file_size,
			file_hash = excluded.file_hash,
			updated_at = CURRENT_TIMESTAMP
	`

	result, err := fis.db.conn.Exec(query, filePath, content, embeddingBytes, model, fileInfo.Size(), fileHash)
	if err != nil {
		return nil, fmt.Errorf("failed to index file: %w", err)
	}

	id, _ := result.LastInsertId()

	return &IndexedFile{
		ID:        int(id),
		FilePath:  filePath,
		Content:   content,
		Embedding: embedding,
		Model:     model,
		FileSize:  fileInfo.Size(),
		FileHash:  fileHash,
		IndexedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (fis *FileIndexStore) GetFile(filePath string) (*IndexedFile, error) {
	query := `
		SELECT id, file_path, content, embedding, model, file_size, file_hash, indexed_at, updated_at
		FROM indexed_files
		WHERE file_path = ?
	`

	var file IndexedFile
	var embeddingBytes []byte

	err := fis.db.conn.QueryRow(query, filePath).Scan(
		&file.ID,
		&file.FilePath,
		&file.Content,
		&embeddingBytes,
		&file.Model,
		&file.FileSize,
		&file.FileHash,
		&file.IndexedAt,
		&file.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("file not indexed: %s", filePath)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get indexed file: %w", err)
	}

	if len(embeddingBytes) > 0 {
		file.Embedding, err = bytesToFloat64Slice(embeddingBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize embedding: %w", err)
		}
	}

	return &file, nil
}

func (fis *FileIndexStore) ListFiles(limit int) ([]IndexedFile, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `
		SELECT id, file_path, content, model, file_size, file_hash, indexed_at, updated_at
		FROM indexed_files
		ORDER BY updated_at DESC
		LIMIT ?
	`

	rows, err := fis.db.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list indexed files: %w", err)
	}
	defer rows.Close()

	var files []IndexedFile
	for rows.Next() {
		var file IndexedFile
		err := rows.Scan(
			&file.ID,
			&file.FilePath,
			&file.Content,
			&file.Model,
			&file.FileSize,
			&file.FileHash,
			&file.IndexedAt,
			&file.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan indexed file: %w", err)
		}
		files = append(files, file)
	}

	return files, rows.Err()
}

func (fis *FileIndexStore) DeleteFile(filePath string) error {
	query := `DELETE FROM indexed_files WHERE file_path = ?`

	result, err := fis.db.conn.Exec(query, filePath)
	if err != nil {
		return fmt.Errorf("failed to delete indexed file: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("file not indexed: %s", filePath)
	}

	return nil
}

func (fis *FileIndexStore) AddChunk(fileID, chunkIndex int, content string, embedding []float64) (*FileChunk, error) {
	embeddingBytes, err := float64SliceToBytes(embedding)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize embedding: %w", err)
	}

	query := `
		INSERT INTO file_chunks (file_id, chunk_index, content, embedding)
		VALUES (?, ?, ?, ?)
	`

	result, err := fis.db.conn.Exec(query, fileID, chunkIndex, content, embeddingBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add file chunk: %w", err)
	}

	id, _ := result.LastInsertId()

	return &FileChunk{
		ID:         int(id),
		FileID:     fileID,
		ChunkIndex: chunkIndex,
		Content:    content,
		Embedding:  embedding,
		CreatedAt:  time.Now(),
	}, nil
}

func (fis *FileIndexStore) GetChunks(fileID int) ([]FileChunk, error) {
	query := `
		SELECT id, file_id, chunk_index, content, embedding, created_at
		FROM file_chunks
		WHERE file_id = ?
		ORDER BY chunk_index ASC
	`

	rows, err := fis.db.conn.Query(query, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file chunks: %w", err)
	}
	defer rows.Close()

	var chunks []FileChunk
	for rows.Next() {
		var chunk FileChunk
		var embeddingBytes []byte

		err := rows.Scan(
			&chunk.ID,
			&chunk.FileID,
			&chunk.ChunkIndex,
			&chunk.Content,
			&embeddingBytes,
			&chunk.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file chunk: %w", err)
		}

		if len(embeddingBytes) > 0 {
			chunk.Embedding, err = bytesToFloat64Slice(embeddingBytes)
			if err != nil {
				return nil, fmt.Errorf("failed to deserialize embedding: %w", err)
			}
		}

		chunks = append(chunks, chunk)
	}

	return chunks, rows.Err()
}

func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
