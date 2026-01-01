// Author: retrozenith <80767544+retrozenith@users.noreply.github.com>
// Description: HTTP handlers for configuration management endpoints

package handlers

import (
	"bluenode-helper/database"
	"encoding/json"
	"log"
	"net/http"
)

type ConfigHandler struct {
	store *database.ConfigStore
}

type SetConfigRequest struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

func NewConfigHandler(store *database.ConfigStore) *ConfigHandler {
	return &ConfigHandler{
		store: store,
	}
}

func (h *ConfigHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	configs, err := h.store.GetAll()
	if err != nil {
		log.Printf("Failed to get all configurations: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, configs)
}

func (h *ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		writeError(w, http.StatusBadRequest, "Configuration key is required")
		return
	}

	config, err := h.store.Get(key)
	if err != nil {
		log.Printf("Failed to get configuration %s: %v", key, err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeSuccess(w, config)
}

func (h *ConfigHandler) Set(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req SetConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Key == "" {
		writeError(w, http.StatusBadRequest, "Configuration key is required")
		return
	}

	if req.Value == "" {
		writeError(w, http.StatusBadRequest, "Configuration value is required")
		return
	}

	if err := h.store.Set(req.Key, req.Value, req.Description); err != nil {
		log.Printf("Failed to set configuration %s: %v", req.Key, err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, map[string]string{
		"status": "success",
		"key":    req.Key,
	})
}

func (h *ConfigHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		writeError(w, http.StatusBadRequest, "Configuration key is required")
		return
	}

	if err := h.store.Delete(key); err != nil {
		log.Printf("Failed to delete configuration %s: %v", key, err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeSuccess(w, map[string]string{
		"status": "deleted",
		"key":    key,
	})
}

func (h *ConfigHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/config", h.GetAll)
	mux.HandleFunc("/config/get", h.Get)
	mux.HandleFunc("/config/set", h.Set)
	mux.HandleFunc("/config/delete", h.Delete)
}
