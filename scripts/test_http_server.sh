#!/bin/bash

# Test script for HTTP-based MCP server
# This script demonstrates how to interact with the HTTP MCP server

SERVER_URL="http://localhost:8080"

echo "=== HTTP MCP Server Test Script ==="
echo "Server URL: $SERVER_URL"
echo ""

# Function to make HTTP requests
make_request() {
    local endpoint="$1"
    local data="$2"
    local method="${3:-POST}"
    
    echo "Making $method request to $endpoint"
    if [ -n "$data" ]; then
        echo "Data: $data"
        echo "Response:"
        curl -s -X "$method" \
             -H "Content-Type: application/json" \
             -d "$data" \
             "$SERVER_URL$endpoint"
    else
        echo "Response:"
        curl -s -X "$method" "$SERVER_URL$endpoint"
    fi
    echo ""
    echo "----------------------------------------"
}

# Test 1: Check server info
echo "1. Testing server info endpoint..."
make_request "/info" "" "GET"

# Test 2: Check health
echo "2. Testing health check..."
make_request "/health" "" "GET"

# Test 3: Check root endpoint
echo "3. Testing root endpoint..."
make_request "/" "" "GET"

# Test 4: Initialize MCP connection
echo "4. Testing MCP initialize..."
initialize_request='{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "http-test-client",
      "version": "1.0.0"
    }
  }
}'
make_request "/mcp" "$initialize_request"

# Test 5: List tools
echo "5. Testing tools/list..."
list_tools_request='{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list",
  "params": {}
}'
make_request "/mcp" "$list_tools_request"

# Test 6: Call echo tool
echo "6. Testing tools/call (echo)..."
echo_tool_request='{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "echo",
    "arguments": {
      "text": "Hello from HTTP client!"
    }
  }
}'
make_request "/mcp" "$echo_tool_request"

# Test 7: List resources
echo "7. Testing resources/list..."
list_resources_request='{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "resources/list",
  "params": {}
}'
make_request "/mcp" "$list_resources_request"

# Test 8: Batch request
echo "8. Testing batch request..."
batch_request='[
  {
    "jsonrpc": "2.0",
    "id": 5,
    "method": "tools/list",
    "params": {}
  },
  {
    "jsonrpc": "2.0",
    "id": 6,
    "method": "resources/list",
    "params": {}
  }
]'
make_request "/mcp" "$batch_request"

echo "=== Test completed ===" 