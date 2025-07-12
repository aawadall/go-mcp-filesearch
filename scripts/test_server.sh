#!/bin/bash

# Test script for the MCP server
# This script sends a sequence of JSON-RPC requests to test the server functionality

echo "Starting MCP Server Test"
echo "========================"

# Build the server if it doesn't exist
if [ ! -f "./mcp-server" ]; then
    echo "Building server..."
    go build -o mcp-server ../cmd/server/main.go
    if [ $? -ne 0 ]; then
        echo "Failed to build server"
        exit 1
    fi
fi

echo "Starting server in background..."
./mcp-server &
SERVER_PID=$!

# Wait a moment for server to start
sleep 1

echo "Sending test requests..."
echo ""

# Create a temporary file to store all requests
TEMP_FILE="/tmp/mcp_requests_$$.json"

# Function to send a request and display response
send_request() {
    local description="$1"
    local request="$2"
    
    echo "=== $description ==="
    echo "Request: $request"
    echo "Response:"
    echo "$request" | timeout 5s ./mcp-server 2>/dev/null || echo "No response received"
    echo ""
}

# Cleanup function
cleanup() {
    echo "Cleaning up..."
    kill $SERVER_PID 2>/dev/null
    rm -f "$TEMP_FILE"
    exit 0
}

# Set up cleanup on script exit
trap cleanup EXIT INT TERM

# Test 1: Initialize
send_request "1. Initialize Connection" '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "resources": {
        "subscribe": true,
        "listChanged": true
      },
      "tools": {
        "listChanged": true
      }
    },
    "clientInfo": {
      "name": "test-client",
      "version": "1.0.0"
    }
  }
}'

# Test 2: List Resources
send_request "2. List Resources" '{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "resources/list",
  "params": {}
}'

# Test 3: Read Resource
send_request "3. Read Resource" '{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "resources/read",
  "params": {
    "uri": "example://test"
  }
}'

# Test 4: List Tools
send_request "4. List Tools" '{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/list",
  "params": {}
}'

# Test 5: Call Echo Tool
send_request "5. Call Echo Tool" '{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "echo",
    "arguments": {
      "text": "Hello, World!"
    }
  }
}'

# Test 6: Call Echo Tool with different message
send_request "6. Call Echo Tool (Different Message)" '{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "tools/call",
  "params": {
    "name": "echo",
    "arguments": {
      "text": "Testing the MCP server!"
    }
  }
}'

# Test 7: Invalid Method (Error Test)
send_request "7. Invalid Method (Error Test)" '{
  "jsonrpc": "2.0",
  "id": 7,
  "method": "invalid/method",
  "params": {}
}'

# Test 8: Invalid Tool (Error Test)
send_request "8. Invalid Tool (Error Test)" '{
  "jsonrpc": "2.0",
  "id": 8,
  "method": "tools/call",
  "params": {
    "name": "nonexistent_tool",
    "arguments": {}
  }
}'

# Test 9: Missing Required Argument (Error Test)
send_request "9. Missing Required Argument (Error Test)" '{
  "jsonrpc": "2.0",
  "id": 9,
  "method": "tools/call",
  "params": {
    "name": "echo",
    "arguments": {}
  }
}'

echo "Test completed!"
echo "===============" 