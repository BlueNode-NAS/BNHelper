// Author: retrozenith <80767544+retrozenith@users.noreply.github.com>
// Description: HTTP handlers for Docker API endpoints

package handlers

import (
	"bluenode-helper/docker"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type DockerHandler struct {
	client *docker.Client
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewDockerHandler() *DockerHandler {
	return &DockerHandler{
		client: docker.NewClient(),
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, APIResponse{
		Success: false,
		Error:   message,
	})
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

func (h *DockerHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if err := h.client.Ping(r.Context()); err != nil {
		log.Printf("Docker ping failed: %v", err)
		writeError(w, http.StatusServiceUnavailable, "Docker daemon is not accessible")
		return
	}

	writeSuccess(w, map[string]string{"status": "ok"})
}

func (h *DockerHandler) Info(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	info, err := h.client.GetInfo(r.Context())
	if err != nil {
		log.Printf("Failed to get Docker info: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, info)
}

func (h *DockerHandler) Version(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	version, err := h.client.GetVersion(r.Context())
	if err != nil {
		log.Printf("Failed to get Docker version: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, version)
}

func (h *DockerHandler) ListContainers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	all := r.URL.Query().Get("all") == "true"
	containers, err := h.client.ListContainers(r.Context(), all)
	if err != nil {
		log.Printf("Failed to list containers: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, containers)
}

func (h *DockerHandler) StartContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	containerID := r.URL.Query().Get("id")
	if containerID == "" {
		writeError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	if err := h.client.StartContainer(r.Context(), containerID); err != nil {
		log.Printf("Failed to start container %s: %v", containerID, err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, map[string]string{"status": "started", "container": containerID})
}

func (h *DockerHandler) StopContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	containerID := r.URL.Query().Get("id")
	if containerID == "" {
		writeError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	timeout := 10
	if t := r.URL.Query().Get("timeout"); t != "" {
		if parsed, err := strconv.Atoi(t); err == nil {
			timeout = parsed
		}
	}

	if err := h.client.StopContainer(r.Context(), containerID, timeout); err != nil {
		log.Printf("Failed to stop container %s: %v", containerID, err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, map[string]string{"status": "stopped", "container": containerID})
}

func (h *DockerHandler) RestartContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	containerID := r.URL.Query().Get("id")
	if containerID == "" {
		writeError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	timeout := 10
	if t := r.URL.Query().Get("timeout"); t != "" {
		if parsed, err := strconv.Atoi(t); err == nil {
			timeout = parsed
		}
	}

	if err := h.client.RestartContainer(r.Context(), containerID, timeout); err != nil {
		log.Printf("Failed to restart container %s: %v", containerID, err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, map[string]string{"status": "restarted", "container": containerID})
}

func (h *DockerHandler) PauseContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	containerID := r.URL.Query().Get("id")
	if containerID == "" {
		writeError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	if err := h.client.PauseContainer(r.Context(), containerID); err != nil {
		log.Printf("Failed to pause container %s: %v", containerID, err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, map[string]string{"status": "paused", "container": containerID})
}

func (h *DockerHandler) UnpauseContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	containerID := r.URL.Query().Get("id")
	if containerID == "" {
		writeError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	if err := h.client.UnpauseContainer(r.Context(), containerID); err != nil {
		log.Printf("Failed to unpause container %s: %v", containerID, err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, map[string]string{"status": "unpaused", "container": containerID})
}

func (h *DockerHandler) RemoveContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	containerID := r.URL.Query().Get("id")
	if containerID == "" {
		writeError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	force := r.URL.Query().Get("force") == "true"
	removeVolumes := r.URL.Query().Get("v") == "true"

	if err := h.client.RemoveContainer(r.Context(), containerID, force, removeVolumes); err != nil {
		log.Printf("Failed to remove container %s: %v", containerID, err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, map[string]string{"status": "removed", "container": containerID})
}

func (h *DockerHandler) ContainerLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	containerID := r.URL.Query().Get("id")
	if containerID == "" {
		writeError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	tail := r.URL.Query().Get("tail")
	if tail == "" {
		tail = "100"
	}
	timestamps := r.URL.Query().Get("timestamps") == "true"

	logs, err := h.client.GetContainerLogs(r.Context(), containerID, tail, timestamps)
	if err != nil {
		log.Printf("Failed to get logs for container %s: %v", containerID, err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, map[string]string{"logs": logs})
}

func (h *DockerHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	images, err := h.client.ListImages(r.Context())
	if err != nil {
		log.Printf("Failed to list images: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, images)
}

func (h *DockerHandler) RemoveImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	imageID := r.URL.Query().Get("id")
	if imageID == "" {
		writeError(w, http.StatusBadRequest, "Image ID is required")
		return
	}

	force := r.URL.Query().Get("force") == "true"

	if err := h.client.RemoveImage(r.Context(), imageID, force); err != nil {
		log.Printf("Failed to remove image %s: %v", imageID, err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, map[string]string{"status": "removed", "image": imageID})
}

func (h *DockerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/docker/ping", h.Ping)
	mux.HandleFunc("/docker/info", h.Info)
	mux.HandleFunc("/docker/version", h.Version)

	mux.HandleFunc("/docker/containers", h.ListContainers)
	mux.HandleFunc("/docker/containers/start", h.StartContainer)
	mux.HandleFunc("/docker/containers/stop", h.StopContainer)
	mux.HandleFunc("/docker/containers/restart", h.RestartContainer)
	mux.HandleFunc("/docker/containers/pause", h.PauseContainer)
	mux.HandleFunc("/docker/containers/unpause", h.UnpauseContainer)
	mux.HandleFunc("/docker/containers/remove", h.RemoveContainer)
	mux.HandleFunc("/docker/containers/logs", h.ContainerLogs)

	mux.HandleFunc("/docker/images", h.ListImages)
	mux.HandleFunc("/docker/images/remove", h.RemoveImage)
}
