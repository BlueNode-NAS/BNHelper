// Author: retrozenith <80767544+retrozenith@users.noreply.github.com>
// Description: HTTP handlers for Ollama AI endpoints

package handlers

import (
	"bluenode-helper/database"
	"bluenode-helper/ollama"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type OllamaHandler struct {
	client         *ollama.Client
	chatStore      *database.ChatStore
	fileIndexStore *database.FileIndexStore
	configStore    *database.ConfigStore
}

type ChatRequest struct {
	Model     string `json:"model"`
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
}

type ChatSessionRequest struct {
	Model string `json:"model"`
	Title string `json:"title,omitempty"`
}

type IndexFileRequest struct {
	FilePath string `json:"file_path"`
	Model    string `json:"model"`
}

func NewOllamaHandler(client *ollama.Client, chatStore *database.ChatStore, fileIndexStore *database.FileIndexStore, configStore *database.ConfigStore) *OllamaHandler {
	return &OllamaHandler{
		client:         client,
		chatStore:      chatStore,
		fileIndexStore: fileIndexStore,
		configStore:    configStore,
	}
}

func (h *OllamaHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if err := h.client.Ping(r.Context()); err != nil {
		log.Printf("Ollama ping failed: %v", err)
		writeError(w, http.StatusServiceUnavailable, "Ollama is not accessible")
		return
	}

	writeSuccess(w, map[string]string{"status": "ok"})
}

func (h *OllamaHandler) ListModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	models, err := h.client.ListModels(r.Context())
	if err != nil {
		log.Printf("Failed to list Ollama models: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, models)
}

func (h *OllamaHandler) Chat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Model == "" {
		config, err := h.configStore.Get("ollama.default_model")
		if err != nil {
			req.Model = "qwen2.5:0.5b"
		} else {
			req.Model = config.Value
		}
	}

	if req.Message == "" {
		writeError(w, http.StatusBadRequest, "Message is required")
		return
	}

	systemPrompt := "You are BlueNode Helper, an AI assistant for the BlueNode Server OS."
	config, err := h.configStore.Get("ollama.system_prompt")
	if err == nil && config.Value != "" {
		systemPrompt = config.Value
	}

	var messages []ollama.Message
	var sessionID string

	if req.SessionID != "" {
		sessionID = req.SessionID
		chatMessages, err := h.chatStore.GetMessages(sessionID)
		if err != nil {
			log.Printf("Failed to get chat messages: %v", err)
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(chatMessages) == 0 {
			messages = append(messages, ollama.Message{
				Role:    "system",
				Content: systemPrompt,
			})
		}

		for _, msg := range chatMessages {
			messages = append(messages, ollama.Message{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	} else {
		session, err := h.chatStore.CreateSession(req.Model, "")
		if err != nil {
			log.Printf("Failed to create chat session: %v", err)
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		sessionID = session.SessionID

		messages = append(messages, ollama.Message{
			Role:    "system",
			Content: systemPrompt,
		})

		h.chatStore.AddMessage(sessionID, "system", systemPrompt)
	}

	messages = append(messages, ollama.Message{
		Role:    "user",
		Content: req.Message,
	})

	h.chatStore.AddMessage(sessionID, "user", req.Message)

	resp, err := h.client.Chat(r.Context(), req.Model, messages)
	if err != nil {
		log.Printf("Failed to chat with Ollama: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.chatStore.AddMessage(sessionID, "assistant", resp.Message.Content)

	writeSuccess(w, map[string]interface{}{
		"session_id": sessionID,
		"message":    resp.Message.Content,
		"model":      resp.Model,
	})
}

func (h *OllamaHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ChatSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Model == "" {
		config, err := h.configStore.Get("ollama.default_model")
		if err != nil {
			req.Model = "qwen2.5:0.5b"
		} else {
			req.Model = config.Value
		}
	}

	session, err := h.chatStore.CreateSession(req.Model, req.Title)
	if err != nil {
		log.Printf("Failed to create chat session: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	systemPrompt := "You are BlueNode Helper, an AI assistant for the BlueNode Server OS."
	config, err := h.configStore.Get("ollama.system_prompt")
	if err == nil && config.Value != "" {
		systemPrompt = config.Value
	}

	h.chatStore.AddMessage(session.SessionID, "system", systemPrompt)

	writeSuccess(w, session)
}

func (h *OllamaHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "Session ID is required")
		return
	}

	session, err := h.chatStore.GetSession(sessionID)
	if err != nil {
		log.Printf("Failed to get chat session: %v", err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	messages, err := h.chatStore.GetMessages(sessionID)
	if err != nil {
		log.Printf("Failed to get chat messages: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, map[string]interface{}{
		"session":  session,
		"messages": messages,
	})
}

func (h *OllamaHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	sessions, err := h.chatStore.ListSessions(limit)
	if err != nil {
		log.Printf("Failed to list chat sessions: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, sessions)
}

func (h *OllamaHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "Session ID is required")
		return
	}

	if err := h.chatStore.DeleteSession(sessionID); err != nil {
		log.Printf("Failed to delete chat session: %v", err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeSuccess(w, map[string]string{
		"status":     "deleted",
		"session_id": sessionID,
	})
}

func (h *OllamaHandler) IndexFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req IndexFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.FilePath == "" {
		writeError(w, http.StatusBadRequest, "File path is required")
		return
	}

	if req.Model == "" {
		req.Model = "nomic-embed-text"
	}

	content, err := os.ReadFile(req.FilePath)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	embedding, err := h.client.GenerateEmbedding(r.Context(), req.Model, string(content))
	if err != nil {
		log.Printf("Failed to generate embedding: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	indexedFile, err := h.fileIndexStore.IndexFile(req.FilePath, string(content), req.Model, embedding)
	if err != nil {
		log.Printf("Failed to index file: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, indexedFile)
}

func (h *OllamaHandler) GetIndexedFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filePath := r.URL.Query().Get("file_path")
	if filePath == "" {
		writeError(w, http.StatusBadRequest, "File path is required")
		return
	}

	file, err := h.fileIndexStore.GetFile(filePath)
	if err != nil {
		log.Printf("Failed to get indexed file: %v", err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeSuccess(w, file)
}

func (h *OllamaHandler) ListIndexedFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	files, err := h.fileIndexStore.ListFiles(limit)
	if err != nil {
		log.Printf("Failed to list indexed files: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, files)
}

func (h *OllamaHandler) DeleteIndexedFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filePath := r.URL.Query().Get("file_path")
	if filePath == "" {
		writeError(w, http.StatusBadRequest, "File path is required")
		return
	}

	if err := h.fileIndexStore.DeleteFile(filePath); err != nil {
		log.Printf("Failed to delete indexed file: %v", err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeSuccess(w, map[string]string{
		"status":    "deleted",
		"file_path": filePath,
	})
}

func (h *OllamaHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ollama/ping", h.Ping)
	mux.HandleFunc("/ollama/models", h.ListModels)
	
	mux.HandleFunc("/ollama/chat", h.Chat)
	mux.HandleFunc("/ollama/sessions", h.ListSessions)
	mux.HandleFunc("/ollama/sessions/create", h.CreateSession)
	mux.HandleFunc("/ollama/sessions/get", h.GetSession)
	mux.HandleFunc("/ollama/sessions/delete", h.DeleteSession)
	
	mux.HandleFunc("/ollama/files/index", h.IndexFile)
	mux.HandleFunc("/ollama/files/get", h.GetIndexedFile)
	mux.HandleFunc("/ollama/files", h.ListIndexedFiles)
	mux.HandleFunc("/ollama/files/delete", h.DeleteIndexedFile)
}
