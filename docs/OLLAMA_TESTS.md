# Quick Ollama Test Commands

## Prerequisites
Make sure Ollama is running and BlueNode Helper service is started.

## Basic Tests

### 1. Ping Ollama
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/ollama/ping
```

### 2. List Models
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/ollama/models
```

---

## Chat Tests

### 3. Start New Chat (Default Model)
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"message":"Hello! What can you help me with?"}' \
  http://localhost/ollama/chat
```

**Save the `session_id` from the response to continue the conversation!**

### 4. Continue Conversation
```bash
# Replace YOUR_SESSION_ID with actual session_id from previous response
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"session_id":"YOUR_SESSION_ID","message":"Tell me more about Docker"}' \
  http://localhost/ollama/chat
```

### 5. List All Sessions
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/ollama/sessions
```

### 6. Get Session with History
```bash
# Replace YOUR_SESSION_ID
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/ollama/sessions/get?session_id=YOUR_SESSION_ID"
```

### 7. Delete Session
```bash
# Replace YOUR_SESSION_ID
curl --unix-socket /var/run/bnhelper.sock \
  -X DELETE \
  "http://localhost/ollama/sessions/delete?session_id=YOUR_SESSION_ID"
```

---

## File Indexing Tests

### 8. Index a File
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"file_path":"/etc/hosts"}' \
  http://localhost/ollama/files/index
```

### 9. List Indexed Files
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/ollama/files
```

### 10. Get Indexed File
```bash
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/ollama/files/get?file_path=/etc/hosts"
```

### 11. Delete Indexed File
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X DELETE \
  "http://localhost/ollama/files/delete?file_path=/etc/hosts"
```

---

## Configuration Tests

### 12. View Current Default Model
```bash
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/config/get?key=ollama.default_model"
```

### 13. Change Default Model
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"key":"ollama.default_model","value":"llama2:latest"}' \
  http://localhost/config/set
```

### 14. View System Prompt
```bash
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/config/get?key=ollama.system_prompt"
```

### 15. Change System Prompt
```bash
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"key":"ollama.system_prompt","value":"You are a helpful Linux server assistant."}' \
  http://localhost/config/set
```

### 16. View All Configurations
```bash
curl --unix-socket /var/run/bnhelper.sock http://localhost/config
```

---

## Full Workflow Example

```bash
# 1. Check if Ollama is accessible
curl --unix-socket /var/run/bnhelper.sock http://localhost/ollama/ping

# 2. List available models
curl --unix-socket /var/run/bnhelper.sock http://localhost/ollama/models

# 3. Start a conversation
RESPONSE=$(curl -s --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"message":"What is Docker?"}' \
  http://localhost/ollama/chat)
echo $RESPONSE

# 4. Extract session_id (requires jq)
SESSION_ID=$(echo $RESPONSE | jq -r '.data.session_id')
echo "Session ID: $SESSION_ID"

# 5. Continue the conversation
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d "{\"session_id\":\"$SESSION_ID\",\"message\":\"How do I list containers?\"}" \
  http://localhost/ollama/chat

# 6. View conversation history
curl --unix-socket /var/run/bnhelper.sock \
  "http://localhost/ollama/sessions/get?session_id=$SESSION_ID" | jq
```

---

## Tips

1. **Pretty print JSON**: Pipe output to `jq` for better readability
   ```bash
   curl --unix-socket /var/run/bnhelper.sock http://localhost/ollama/sessions | jq
   ```

2. **Save responses**: Use `-o` to save responses to files
   ```bash
   curl --unix-socket /var/run/bnhelper.sock http://localhost/ollama/models -o models.json
   ```

3. **Verbose output**: Add `-v` for debugging
   ```bash
   curl -v --unix-socket /var/run/bnhelper.sock http://localhost/ollama/ping
   ```

4. **View the test script**: For all available commands
   ```bash
   bash /root/BNHelper/docs/ollama_test_commands.sh
   ```
