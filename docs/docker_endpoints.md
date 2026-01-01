# Docker API Endpoints

BlueNode Helper provides a RESTful API over Unix domain socket for Docker management.

## Base Information

- **Socket Path**: `/var/run/bnhelper.sock`
- **Protocol**: HTTP over Unix socket
- **Response Format**: JSON

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

## General Endpoints

### Health Check

Check if the service is running.

**Endpoint**: `GET /health`

**Response**:
```
OK
```

---

### Root

API information endpoint.

**Endpoint**: `GET /`

**Response**:
```
BlueNode Helper API
```

---

## Docker System Endpoints

### Ping Docker Daemon

Check if Docker daemon is accessible.

**Endpoint**: `GET /docker/ping`

**Example Request**:
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/docker/ping
```

**Example Response**:
```json
{
  "success": true,
  "data": {
    "status": "ok"
  }
}
```

---

### Get Docker Version

Retrieve Docker daemon version information.

**Endpoint**: `GET /docker/version`

**Example Request**:
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/docker/version
```

**Example Response**:
```json
{
  "success": true,
  "data": {
    "Version": "24.0.5",
    "ApiVersion": "1.43",
    "GitCommit": "ced0996",
    "GoVersion": "go1.21.0",
    "Os": "linux",
    "Arch": "amd64",
    "KernelVersion": "6.5.0",
    "BuildTime": "2023-07-20T10:15:30.000000000+00:00"
  }
}
```

---

### Get Docker System Info

Retrieve detailed Docker system information.

**Endpoint**: `GET /docker/info`

**Example Request**:
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/docker/info
```

**Example Response**:
```json
{
  "success": true,
  "data": {
    "ID": "...",
    "Containers": 5,
    "ContainersRunning": 2,
    "ContainersPaused": 0,
    "ContainersStopped": 3,
    "Images": 15,
    "Driver": "overlay2",
    "MemoryLimit": true,
    "SwapLimit": true,
    "KernelVersion": "6.5.0",
    "OperatingSystem": "Fedora Linux",
    "Architecture": "x86_64",
    "NCPU": 8,
    "MemTotal": 16777216000
  }
}
```

---

## Container Endpoints

### List Containers

List Docker containers.

**Endpoint**: `GET /docker/containers`

**Query Parameters**:
- `all` (boolean, optional): Show all containers (default: false, only running containers)

**Example Request**:
```bash
# List running containers only
curl --unix-socket /var/run/bnhelper.sock http://localhost/docker/containers

# List all containers (including stopped)
curl --unix-socket /var/run/bnhelper.sock "http://localhost/docker/containers?all=true"
```

**Example Response**:
```json
{
  "success": true,
  "data": [
    {
      "Id": "325e144ed56801c8520bc82c06e0d0661faa939a18ed2ada11244163a8012ed7",
      "Names": ["/iamfoo"],
      "Image": "containous/whoami",
      "ImageID": "sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
      "Command": "/whoami",
      "Created": 1767283526,
      "State": "running",
      "Status": "Up 14 minutes",
      "Ports": [
        {
          "IP": "0.0.0.0",
          "PrivatePort": 80,
          "PublicPort": 32768,
          "Type": "tcp"
        }
      ]
    }
  ]
}
```

---

### Start Container

Start a stopped container.

**Endpoint**: `POST /docker/containers/start`

**Query Parameters**:
- `id` (string, required): Container ID or name

**Example Request**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  "http://localhost/docker/containers/start?id=325e144ed568"
```

**Example Response**:
```json
{
  "success": true,
  "data": {
    "status": "started",
    "container": "325e144ed568"
  }
}
```

---

### Stop Container

Stop a running container.

**Endpoint**: `POST /docker/containers/stop`

**Query Parameters**:
- `id` (string, required): Container ID or name
- `timeout` (integer, optional): Timeout in seconds before killing (default: 10)

**Example Request**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  "http://localhost/docker/containers/stop?id=325e144ed568&timeout=30"
```

**Example Response**:
```json
{
  "success": true,
  "data": {
    "status": "stopped",
    "container": "325e144ed568"
  }
}
```

---

### Restart Container

Restart a container.

**Endpoint**: `POST /docker/containers/restart`

**Query Parameters**:
- `id` (string, required): Container ID or name
- `timeout` (integer, optional): Timeout in seconds before killing (default: 10)

**Example Request**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  "http://localhost/docker/containers/restart?id=325e144ed568"
```

**Example Response**:
```json
{
  "success": true,
  "data": {
    "status": "restarted",
    "container": "325e144ed568"
  }
}
```

---

### Pause Container

Pause a running container.

**Endpoint**: `POST /docker/containers/pause`

**Query Parameters**:
- `id` (string, required): Container ID or name

**Example Request**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  "http://localhost/docker/containers/pause?id=325e144ed568"
```

**Example Response**:
```json
{
  "success": true,
  "data": {
    "status": "paused",
    "container": "325e144ed568"
  }
}
```

---

### Unpause Container

Unpause a paused container.

**Endpoint**: `POST /docker/containers/unpause`

**Query Parameters**:
- `id` (string, required): Container ID or name

**Example Request**:
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  "http://localhost/docker/containers/unpause?id=325e144ed568"
```

**Example Response**:
```json
{
  "success": true,
  "data": {
    "status": "unpaused",
    "container": "325e144ed568"
  }
}
```

---

### Remove Container

Remove a container.

**Endpoint**: `DELETE /docker/containers/remove`

**Query Parameters**:
- `id` (string, required): Container ID or name
- `force` (boolean, optional): Force removal of running container (default: false)
- `v` (boolean, optional): Remove associated volumes (default: false)

**Example Request**:
```bash
# Remove stopped container
curl --unix-socket /var/run/bnhelper.sock \
  -X DELETE \
  "http://localhost/docker/containers/remove?id=325e144ed568"

# Force remove running container with volumes
curl --unix-socket /var/run/bnhelper.sock \
  -X DELETE \
  "http://localhost/docker/containers/remove?id=325e144ed568&force=true&v=true"
```

**Example Response**:
```json
{
  "success": true,
  "data": {
    "status": "removed",
    "container": "325e144ed568"
  }
}
```

---

### Get Container Logs

Retrieve logs from a container.

**Endpoint**: `GET /docker/containers/logs`

**Query Parameters**:
- `id` (string, required): Container ID or name
- `tail` (string, optional): Number of lines to show from the end (default: "100")
- `timestamps` (boolean, optional): Show timestamps (default: false)

**Example Request**:
```bash
# Get last 100 lines
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/docker/containers/logs?id=325e144ed568"

# Get last 50 lines with timestamps
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/docker/containers/logs?id=325e144ed568&tail=50&timestamps=true"
```

**Example Response**:
```json
{
  "success": true,
  "data": {
    "logs": "2026-01-01T16:20:15.123456789Z Starting server...\n2026-01-01T16:20:15.234567890Z Server listening on port 80\n"
  }
}
```

---

## Image Endpoints

### List Images

List Docker images.

**Endpoint**: `GET /docker/images`

**Example Request**:
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/docker/images
```

**Example Response**:
```json
{
  "success": true,
  "data": [
    {
      "Id": "sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
      "ParentId": "",
      "RepoTags": ["containous/whoami:latest"],
      "RepoDigests": ["containous/whoami@sha256:abc123..."],
      "Created": 1767283500,
      "Size": 6710886,
      "VirtualSize": 6710886,
      "SharedSize": 0,
      "Labels": {},
      "Containers": 1
    }
  ]
}
```

---

### Remove Image

Remove a Docker image.

**Endpoint**: `DELETE /docker/images/remove`

**Query Parameters**:
- `id` (string, required): Image ID or name
- `force` (boolean, optional): Force removal (default: false)

**Example Request**:
```bash
# Remove unused image
curl --unix-socket /var/run/bnhelper.sock \
  -X DELETE \
  "http://localhost/docker/images/remove?id=containous/whoami"

# Force remove image
curl --unix-socket /var/run/bnhelper.sock \
  -X DELETE \
  "http://localhost/docker/images/remove?id=7d6a3c8f9147&force=true"
```

**Example Response**:
```json
{
  "success": true,
  "data": {
    "status": "removed",
    "image": "containous/whoami"
  }
}
```

---

## Error Responses

When an error occurs, the API returns an error response:

```json
{
  "success": false,
  "error": "Container ID is required"
}
```

### Common HTTP Status Codes

- `200 OK`: Request successful
- `400 Bad Request`: Invalid parameters
- `405 Method Not Allowed`: Wrong HTTP method used
- `500 Internal Server Error`: Docker daemon error or internal error
- `503 Service Unavailable`: Docker daemon not accessible

---

## Using with curl

### Basic Request
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/docker/ping
```

### Pretty Print with jq
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/docker/containers | jq .
```

### Set Socket Variable
```bash
SOCKET=/var/run/bnhelper.sock
curl --unix-socket $SOCKET http://localhost/docker/version | jq .
```

---

## Example Workflow

```bash
# Set socket path
SOCKET=/var/run/bnhelper.sock

# 1. Check Docker is accessible
curl --unix-socket $SOCKET http://localhost/docker/ping | jq .

# 2. List all containers
curl --unix-socket $SOCKET "http://localhost/docker/containers?all=true" | jq .

# 3. Start a container
curl --unix-socket $SOCKET -X POST \
  "http://localhost/docker/containers/start?id=mycontainer" | jq .

# 4. Check logs
curl --unix-socket $SOCKET \
  "http://localhost/docker/containers/logs?id=mycontainer&tail=50" | jq .

# 5. Stop the container
curl --unix-socket $SOCKET -X POST \
  "http://localhost/docker/containers/stop?id=mycontainer" | jq .

# 6. List images
curl --unix-socket $SOCKET http://localhost/docker/images | jq .
```

---

## Security Notes

- The Unix socket `/var/run/bnhelper.sock` has permissions `0660` (owner and group only)
- The service runs as root to access Docker daemon
- Ensure proper access control to the socket file
- Consider using groups to manage access to the socket

---

## Troubleshooting

### Socket not found
```bash
# Check if service is running
sudo systemctl status bluenode-helper

# Check socket exists
ls -la /var/run/bnhelper.sock
```

### Permission denied
```bash
# Check socket permissions
ls -la /var/run/bnhelper.sock

# Run with sudo if needed
sudo curl --unix-socket /var/run/bnhelper.sock http://localhost/docker/ping
```

### Docker daemon not accessible
```bash
# Check Docker daemon is running
sudo systemctl status docker

# Test direct Docker access
docker ps
```
