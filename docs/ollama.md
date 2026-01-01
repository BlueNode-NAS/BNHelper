# Ollama API Endpoints

BlueNode Helper provides AI capabilities through Ollama integration with chat sessions and file indexing.

## Base Information

- **Socket Path**: `/var/run/bnhelper.sock`
- **Protocol**: HTTP over Unix socket
- **Response Format**: JSON
- **AI Database Location**: `/var/lib/bnhelper/ai.db`
- **Ollama URL**: `http://localhost:11434` (default)
- **Default Model**: `qwen2.5:0.5b` (configurable via database)
- **System Prompt**: Configurable via database configuration

## Configuration

The Ollama integration uses configuration values from the database:

- `ollama.default_model`: Default model for chat (default: "qwen2.5:0.5b")
- `ollama.system_prompt`: System prompt for AI assistant (default: "You are BlueNode Helper, an AI assistant for the BlueNode Server OS.")

These can be changed using the configuration endpoints (see database.md).

**Example - Change default model**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"key":"ollama.default_model","value":"llama2:latest","description":"Default Ollama model for chat"}' \
  http://localhost/config/set
```

**Example - Change system prompt**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"key":"ollama.system_prompt","value":"You are a helpful assistant for server management.","description":"System prompt for Ollama chat"}' \
  http://localhost/config/set
```

## Response Schema

All endpoints return responses in the following format:

```json
{
  "success": true,
  "data": {},
  "error": ""
}
```

---

## Ollama Endpoints

### Ping Ollama

Check if Ollama service is accessible.

**Endpoint**: `GET /ollama/ping`

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "ok"
  }
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/ollama/ping
```

---

### List Models

Get available Ollama models.

**Endpoint**: `GET /ollama/models`

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "name": "llama2:latest",
      "modified_at": "2026-01-01T10:00:00Z",
      "size": 3826793677
    }
  ]
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/ollama/models
```

---

## Chat Endpoints

### Chat with AI

Send a message and get a response. Automatically manages session history.

**Endpoint**: `POST /ollama/chat`

**Request Body**:
```json
{
  "model": "llama2:latest",
  "message": "Hello, how are you?",
  "session_id": "optional-session-id"
}
```

**Fields**:
- `model` (optional): Ollama model name (defaults to configured `ollama.default_model`)
- `message` (required): User message
- `session_id` (optional): Existing session ID to continue conversation

**Notes**:
- If `model` is not provided, uses the default model from database configuration (`ollama.default_model`)
- System prompt is automatically included from database configuration (`ollama.system_prompt`)
- System prompt is only added at the start of new conversations

**Response**:
```json
{
  "success": true,
  "data": {
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "I'm doing well, thank you for asking!",
    "model": "llama2:latest"
  }
}
```

**Example**:
```bash
# Start new conversation
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"model":"llama2:latest","message":"Hello!"}' \
  http://localhost/ollama/chat

# Continue conversation
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"model":"llama2:latest","message":"Tell me a joke","session_id":"550e8400-e29b-41d4-a716-446655440000"}' \
  http://localhost/ollama/chat
```

---

### Create Chat Session

Manually create a new chat session.

**Endpoint**: `POST /ollama/sessions/create`

**Request Body**:
```json
{
  "model": "llama2:latest",
  "title": "My Conversation"
}
```

**Fields**:
- `model` (optional): Ollama model name (defaults to configured `ollama.default_model`)
- `title` (optional): Human-readable session title

**Notes**:
- If `model` is not provided, uses the default model from database configuration
- System prompt is automatically added from database configuration

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "model": "llama2:latest",
    "title": "My Conversation",
    "created_at": "2026-01-01T18:00:00Z",
    "updated_at": "2026-01-01T18:00:00Z"
  }
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"model":"llama2:latest","title":"Support Chat"}' \
  http://localhost/ollama/sessions/create
```

---

### Get Chat Session

Retrieve a session with all its messages.

**Endpoint**: `GET /ollama/sessions/get?session_id={session_id}`

**Query Parameters**:
- `session_id` (required): Session ID to retrieve

**Response**:
```json
{
  "success": true,
  "data": {
    "session": {
      "id": 1,
      "session_id": "550e8400-e29b-41d4-a716-446655440000",
      "model": "llama2:latest",
      "title": "My Conversation",
      "created_at": "2026-01-01T18:00:00Z",
      "updated_at": "2026-01-01T18:05:00Z"
    },
    "messages": [
      {
        "id": 1,
        "session_id": "550e8400-e29b-41d4-a716-446655440000",
        "role": "user",
        "content": "Hello!",
        "created_at": "2026-01-01T18:00:00Z"
      },
      {
        "id": 2,
        "session_id": "550e8400-e29b-41d4-a716-446655440000",
        "role": "assistant",
        "content": "Hi there! How can I help you?",
        "created_at": "2026-01-01T18:00:01Z"
      }
    ]
  }
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/ollama/sessions/get?session_id=550e8400-e29b-41d4-a716-446655440000"
```

---

### List Chat Sessions

Get all chat sessions ordered by most recent.

**Endpoint**: `GET /ollama/sessions?limit={limit}`

**Query Parameters**:
- `limit` (optional): Maximum number of sessions (default: 50)

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 2,
      "session_id": "660e8400-e29b-41d4-a716-446655440001",
      "model": "llama2:latest",
      "title": "Recent Chat",
      "created_at": "2026-01-01T19:00:00Z",
      "updated_at": "2026-01-01T19:10:00Z"
    },
    {
      "id": 1,
      "session_id": "550e8400-e29b-41d4-a716-446655440000",
      "model": "llama2:latest",
      "title": "Older Chat",
      "created_at": "2026-01-01T18:00:00Z",
      "updated_at": "2026-01-01T18:05:00Z"
    }
  ]
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/ollama/sessions?limit=10"
```

---

### Delete Chat Session

Delete a session and all its messages.

**Endpoint**: `DELETE /ollama/sessions/delete?session_id={session_id}`

**Query Parameters**:
- `session_id` (required): Session ID to delete

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "deleted",
    "session_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X DELETE \
  "http://localhost/ollama/sessions/delete?session_id=550e8400-e29b-41d4-a716-446655440000"
```

---

## File Indexing Endpoints

### Index File

Index a file with embeddings for semantic search.

**Endpoint**: `POST /ollama/files/index`

**Request Body**:
```json
{
  "file_path": "/path/to/document.txt",
  "model": "nomic-embed-text"
}
```

**Fields**:
- `file_path` (required): Absolute path to the file
- `model` (optional): Embedding model (default: "nomic-embed-text")

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "file_path": "/path/to/document.txt",
    "content": "File content here...",
    "model": "nomic-embed-text",
    "file_size": 1024,
    "file_hash": "a3b2c1d4e5f6...",
    "indexed_at": "2026-01-01T18:00:00Z",
    "updated_at": "2026-01-01T18:00:00Z"
  }
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"file_path":"/etc/hosts","model":"nomic-embed-text"}' \
  http://localhost/ollama/files/index
```

---

### Get Indexed File

Retrieve an indexed file with its embedding.

**Endpoint**: `GET /ollama/files/get?file_path={file_path}`

**Query Parameters**:
- `file_path` (required): Path of the indexed file

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "file_path": "/path/to/document.txt",
    "content": "File content...",
    "embedding": [0.1, 0.2, 0.3, ...],
    "model": "nomic-embed-text",
    "file_size": 1024,
    "file_hash": "a3b2c1d4e5f6...",
    "indexed_at": "2026-01-01T18:00:00Z",
    "updated_at": "2026-01-01T18:00:00Z"
  }
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/ollama/files/get?file_path=/etc/hosts"
```

---

### List Indexed Files

Get all indexed files.

**Endpoint**: `GET /ollama/files?limit={limit}`

**Query Parameters**:
- `limit` (optional): Maximum number of files (default: 100)

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 2,
      "file_path": "/path/to/recent.txt",
      "content": "Content...",
      "model": "nomic-embed-text",
      "file_size": 2048,
      "file_hash": "b4c3d2e1f0a9...",
      "indexed_at": "2026-01-01T19:00:00Z",
      "updated_at": "2026-01-01T19:00:00Z"
    }
  ]
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/ollama/files?limit=50"
```

---

### Delete Indexed File

Remove a file from the index.

**Endpoint**: `DELETE /ollama/files/delete?file_path={file_path}`

**Query Parameters**:
- `file_path` (required): Path of the file to delete

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "deleted",
    "file_path": "/path/to/document.txt"
  }
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X DELETE \
  "http://localhost/ollama/files/delete?file_path=/etc/hosts"
```

---

## Database Schema

### chat_sessions Table

| Column     | Type     | Description                       |
|------------|----------|-----------------------------------|
| id         | INTEGER  | Auto-incrementing primary key     |
| session_id | TEXT     | Unique session identifier (UUID)  |
| model      | TEXT     | Ollama model name                 |
| title      | TEXT     | Optional session title            |
| created_at | DATETIME | Creation timestamp                |
| updated_at | DATETIME | Last update timestamp             |

### chat_messages Table

| Column     | Type     | Description                       |
|------------|----------|-----------------------------------|
| id         | INTEGER  | Auto-incrementing primary key     |
| session_id | TEXT     | Foreign key to chat_sessions      |
| role       | TEXT     | "system", "user", or "assistant"  |
| content    | TEXT     | Message content                   |
| created_at | DATETIME | Creation timestamp                |

### indexed_files Table

| Column     | Type     | Description                       |
|------------|----------|-----------------------------------|
| id         | INTEGER  | Auto-incrementing primary key     |
| file_path  | TEXT     | Unique file path                  |
| content    | TEXT     | File content                      |
| embedding  | BLOB     | Serialized embedding vector       |
| model      | TEXT     | Embedding model used              |
| file_size  | INTEGER  | File size in bytes                |
| file_hash  | TEXT     | SHA256 hash of file               |
| indexed_at | DATETIME | Index creation timestamp          |
| updated_at | DATETIME | Last update timestamp             |

### file_chunks Table

| Column      | Type     | Description                       |
|-------------|----------|-----------------------------------|
| id          | INTEGER  | Auto-incrementing primary key     |
| file_id     | INTEGER  | Foreign key to indexed_files      |
| chunk_index | INTEGER  | Chunk position in file            |
| content     | TEXT     | Chunk content                     |
| embedding   | BLOB     | Serialized embedding vector       |
| created_at  | DATETIME | Creation timestamp                |

---

## Error Handling

Common HTTP status codes:

- `200 OK`: Request successful
- `400 Bad Request`: Missing or invalid parameters
- `404 Not Found`: Resource not found
- `405 Method Not Allowed`: Invalid HTTP method
- `500 Internal Server Error`: Server or Ollama error
- `503 Service Unavailable`: Ollama is not accessible

All errors include a descriptive message in the `error` field.
