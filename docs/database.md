# Database API Endpoints

BlueNode Helper provides a configuration storage API over Unix domain socket using SQLite database.

## Base Information

- **Socket Path**: `/var/run/bnhelper.sock`
- **Protocol**: HTTP over Unix socket
- **Response Format**: JSON
- **Database Location**: `/var/lib/bnhelper/bnhelper.db`

## Response Schema

All endpoints return responses in the following format:

```json
{
  "success": true,
  "data": {},
  "error": ""
}
```

- `success` (boolean): Indicates if the request was successful
- `data` (object/array): Response data (only present on success)
- `error` (string): Error message (only present on failure)

---

## Configuration Endpoints

### Get All Configurations

Retrieve all stored configurations.

**Endpoint**: `GET /config`

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "key": "app.name",
      "value": "BlueNode Helper",
      "description": "Application name",
      "created_at": "2026-01-01T18:00:00Z",
      "updated_at": "2026-01-01T18:00:00Z"
    },
    {
      "id": 2,
      "key": "app.version",
      "value": "1.0.0",
      "description": "Application version",
      "created_at": "2026-01-01T18:05:00Z",
      "updated_at": "2026-01-01T18:10:00Z"
    }
  ]
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/config
```

---

### Get Single Configuration

Retrieve a specific configuration by key.

**Endpoint**: `GET /config/get?key={key}`

**Query Parameters**:
- `key` (required): Configuration key to retrieve

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "key": "app.name",
    "value": "BlueNode Helper",
    "description": "Application name",
    "created_at": "2026-01-01T18:00:00Z",
    "updated_at": "2026-01-01T18:00:00Z"
  }
}
```

**Error Response** (404 Not Found):
```json
{
  "success": false,
  "error": "configuration not found: app.invalid"
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/config/get?key=app.name"
```

---

### Set Configuration

Create or update a configuration. If the key already exists, it will be updated with the new value.

**Endpoint**: `POST /config/set` or `PUT /config/set`

**Request Body**:
```json
{
  "key": "app.name",
  "value": "BlueNode Helper",
  "description": "Application name"
}
```

**Fields**:
- `key` (required): Configuration key (must be unique)
- `value` (required): Configuration value
- `description` (optional): Human-readable description

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "success",
    "key": "app.name"
  }
}
```

**Error Response** (400 Bad Request):
```json
{
  "success": false,
  "error": "Configuration key is required"
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"key":"app.name","value":"BlueNode Helper","description":"Application name"}' \
  http://localhost/config/set
```

---

### Delete Configuration

Delete a configuration by key.

**Endpoint**: `DELETE /config/delete?key={key}`

**Query Parameters**:
- `key` (required): Configuration key to delete

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "deleted",
    "key": "app.name"
  }
}
```

**Error Response** (404 Not Found):
```json
{
  "success": false,
  "error": "configuration not found: app.invalid"
}
```

**Example**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X DELETE \
  "http://localhost/config/delete?key=app.name"
```

---

## Database Schema

### configurations Table

The configurations table stores key-value pairs with metadata.

| Column      | Type     | Description                          |
|-------------|----------|--------------------------------------|
| id          | INTEGER  | Auto-incrementing primary key        |
| key         | TEXT     | Unique configuration key             |
| value       | TEXT     | Configuration value                  |
| description | TEXT     | Optional description                 |
| created_at  | DATETIME | Timestamp when created (auto-set)    |
| updated_at  | DATETIME | Timestamp when updated (auto-update) |

**Indexes**:
- Primary key on `id`
- Unique index on `key`

**Triggers**:
- `update_configurations_timestamp`: Automatically updates `updated_at` on row updates

---

## Usage Examples

### Store Multiple Configurations

```bash
# Set database host
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"key":"db.host","value":"localhost","description":"Database host"}' \
  http://localhost/config/set

# Set database port
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"key":"db.port","value":"5432","description":"Database port"}' \
  http://localhost/config/set

# List all configurations
curl --unix-socket /var/run/bnhelper.sock http://localhost/config
```

### Update Existing Configuration

```bash
# Update will replace the existing value and update the timestamp
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"key":"db.port","value":"5433","description":"Database port (updated)"}' \
  http://localhost/config/set
```

### Retrieve and Delete

```bash
# Get specific configuration
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/config/get?key=db.host"

# Delete configuration
curl --unix-socket /var/run/bnhelper.sock \
  -X DELETE \
  "http://localhost/config/delete?key=db.host"
```

---

## Error Handling

Common HTTP status codes:

- `200 OK`: Request successful
- `400 Bad Request`: Missing or invalid parameters
- `404 Not Found`: Configuration key not found
- `405 Method Not Allowed`: Invalid HTTP method
- `500 Internal Server Error`: Database or server error

All errors include a descriptive message in the `error` field of the response.
