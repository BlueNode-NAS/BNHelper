#!/bin/bash
# Ollama API Test Commands for BlueNode Helper

SOCKET="/var/run/bnhelper.sock"

echo "=== Ollama API Tests ==="
echo ""

# Test 1: Ping Ollama
echo "1. Ping Ollama service"
echo "Command:"
echo "curl --unix-socket $SOCKET http://localhost/ollama/ping"
echo ""

# Test 2: List Models
echo "2. List available Ollama models"
echo "Command:"
echo "curl --unix-socket $SOCKET http://localhost/ollama/models"
echo ""

# Test 3: Chat - Start new conversation (uses default model)
echo "3. Start new chat conversation (default model)"
echo "Command:"
cat << 'EOF'
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"message":"Hello! What can you help me with?"}' \
  http://localhost/ollama/chat
EOF
echo ""
echo ""

# Test 4: Chat - With specific model
echo "4. Start chat with specific model"
echo "Command:"
cat << 'EOF'
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"model":"qwen2.5:0.5b","message":"Explain what BlueNode is"}' \
  http://localhost/ollama/chat
EOF
echo ""
echo ""

# Test 5: Continue conversation
echo "5. Continue existing conversation (replace SESSION_ID)"
echo "Command:"
cat << 'EOF'
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"session_id":"YOUR_SESSION_ID_HERE","message":"Tell me more"}' \
  http://localhost/ollama/chat
EOF
echo ""
echo ""

# Test 6: Create session manually
echo "6. Create chat session manually"
echo "Command:"
cat << 'EOF'
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Session"}' \
  http://localhost/ollama/sessions/create
EOF
echo ""
echo ""

# Test 7: List sessions
echo "7. List all chat sessions"
echo "Command:"
echo "curl --unix-socket $SOCKET 'http://localhost/ollama/sessions?limit=10'"
echo ""

# Test 8: Get session details
echo "8. Get session with messages (replace SESSION_ID)"
echo "Command:"
echo "curl --unix-socket $SOCKET 'http://localhost/ollama/sessions/get?session_id=YOUR_SESSION_ID_HERE'"
echo ""

# Test 9: Delete session
echo "9. Delete chat session (replace SESSION_ID)"
echo "Command:"
echo "curl --unix-socket $SOCKET -X DELETE 'http://localhost/ollama/sessions/delete?session_id=YOUR_SESSION_ID_HERE'"
echo ""

# Test 10: Index a file
echo "10. Index a file for semantic search"
echo "Command:"
cat << 'EOF'
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"file_path":"/etc/hosts","model":"nomic-embed-text"}' \
  http://localhost/ollama/files/index
EOF
echo ""
echo ""

# Test 11: Get indexed file
echo "11. Get indexed file details"
echo "Command:"
echo "curl --unix-socket $SOCKET 'http://localhost/ollama/files/get?file_path=/etc/hosts'"
echo ""

# Test 12: List indexed files
echo "12. List all indexed files"
echo "Command:"
echo "curl --unix-socket $SOCKET 'http://localhost/ollama/files?limit=20'"
echo ""

# Test 13: Delete indexed file
echo "13. Delete indexed file"
echo "Command:"
echo "curl --unix-socket $SOCKET -X DELETE 'http://localhost/ollama/files/delete?file_path=/etc/hosts'"
echo ""

# Configuration Tests
echo "=== Configuration Tests ==="
echo ""

# Test 14: Get default model
echo "14. Get current default model"
echo "Command:"
echo "curl --unix-socket $SOCKET 'http://localhost/config/get?key=ollama.default_model'"
echo ""

# Test 15: Change default model
echo "15. Change default model to llama2"
echo "Command:"
cat << 'EOF'
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"key":"ollama.default_model","value":"llama2:latest","description":"Default Ollama model"}' \
  http://localhost/config/set
EOF
echo ""
echo ""

# Test 16: Get system prompt
echo "16. Get current system prompt"
echo "Command:"
echo "curl --unix-socket $SOCKET 'http://localhost/config/get?key=ollama.system_prompt'"
echo ""

# Test 17: Change system prompt
echo "17. Change system prompt"
echo "Command:"
cat << 'EOF'
curl --unix-socket /var/run/bnhelper.sock \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"key":"ollama.system_prompt","value":"You are a helpful server management assistant.","description":"Custom system prompt"}' \
  http://localhost/config/set
EOF
echo ""
echo ""

# Test 18: View all configurations
echo "18. View all configurations"
echo "Command:"
echo "curl --unix-socket $SOCKET http://localhost/config"
echo ""

echo "=== End of Test Commands ==="
