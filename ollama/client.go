// Author: retrozenith <80767544+retrozenith@users.noreply.github.com>
// Description: Ollama client for AI operations

package ollama

import (
"bytes"
"context"
"encoding/json"
"fmt"
"io"
"net/http"
"time"
)

const (
defaultOllamaURL = "http://localhost:11434"
)

type Client struct {
baseURL    string
httpClient *http.Client
}

type ChatRequest struct {
Model    string    `json:"model"`
Messages []Message `json:"messages"`
Stream   bool      `json:"stream"`
}

type Message struct {
Role    string `json:"role"`
Content string `json:"content"`
}

type ChatResponse struct {
Model     string    `json:"model"`
CreatedAt time.Time `json:"created_at"`
Message   Message   `json:"message"`
Done      bool      `json:"done"`
}

type EmbeddingRequest struct {
Model  string `json:"model"`
Prompt string `json:"prompt"`
}

type EmbeddingResponse struct {
Embedding []float64 `json:"embedding"`
}

type ModelListResponse struct {
Models []Model `json:"models"`
}

type Model struct {
Name       string    `json:"name"`
ModifiedAt time.Time `json:"modified_at"`
Size       int64     `json:"size"`
}

func NewClient(baseURL string) *Client {
if baseURL == "" {
baseURL = defaultOllamaURL
}

return &Client{
baseURL: baseURL,
httpClient: &http.Client{
Timeout: 5 * time.Minute,
},
}
}

func (c *Client) Chat(ctx context.Context, model string, messages []Message) (*ChatResponse, error) {
req := ChatRequest{
Model:    model,
Messages: messages,
Stream:   false,
}

body, err := json.Marshal(req)
if err != nil {
return nil, fmt.Errorf("failed to marshal request: %w", err)
}

httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/chat", bytes.NewReader(body))
if err != nil {
return nil, fmt.Errorf("failed to create request: %w", err)
}

httpReq.Header.Set("Content-Type", "application/json")

resp, err := c.httpClient.Do(httpReq)
if err != nil {
return nil, fmt.Errorf("failed to send request: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
bodyBytes, _ := io.ReadAll(resp.Body)
return nil, fmt.Errorf("ollama request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
}

var chatResp ChatResponse
if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
return nil, fmt.Errorf("failed to decode response: %w", err)
}

return &chatResp, nil
}

func (c *Client) GenerateEmbedding(ctx context.Context, model, text string) ([]float64, error) {
req := EmbeddingRequest{
Model:  model,
Prompt: text,
}

body, err := json.Marshal(req)
if err != nil {
return nil, fmt.Errorf("failed to marshal request: %w", err)
}

httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/embeddings", bytes.NewReader(body))
if err != nil {
return nil, fmt.Errorf("failed to create request: %w", err)
}

httpReq.Header.Set("Content-Type", "application/json")

resp, err := c.httpClient.Do(httpReq)
if err != nil {
return nil, fmt.Errorf("failed to send request: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
bodyBytes, _ := io.ReadAll(resp.Body)
return nil, fmt.Errorf("ollama embedding request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
}

var embResp EmbeddingResponse
if err := json.NewDecoder(resp.Body).Decode(&embResp); err != nil {
return nil, fmt.Errorf("failed to decode response: %w", err)
}

return embResp.Embedding, nil
}

func (c *Client) ListModels(ctx context.Context) ([]Model, error) {
httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
if err != nil {
return nil, fmt.Errorf("failed to create request: %w", err)
}

resp, err := c.httpClient.Do(httpReq)
if err != nil {
return nil, fmt.Errorf("failed to send request: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
bodyBytes, _ := io.ReadAll(resp.Body)
return nil, fmt.Errorf("ollama list models request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
}

var listResp ModelListResponse
if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
return nil, fmt.Errorf("failed to decode response: %w", err)
}

return listResp.Models, nil
}

func (c *Client) Ping(ctx context.Context) error {
httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
if err != nil {
return fmt.Errorf("failed to create request: %w", err)
}

resp, err := c.httpClient.Do(httpReq)
if err != nil {
return fmt.Errorf("ollama is not accessible: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
return fmt.Errorf("ollama returned status %d", resp.StatusCode)
}

return nil
}
