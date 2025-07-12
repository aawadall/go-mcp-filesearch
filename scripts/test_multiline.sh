#!/bin/bash

echo "Testing MCP Server Multiline and Batch Parsing"
echo "=============================================="

# Build the server if it doesn't exist
if [ ! -f "./mcp-server" ]; then
    echo "Building server..."
    go build -o mcp-server ../cmd/server/main.go
    if [ $? -ne 0 ]; then
        echo "Failed to build server"
        exit 1
    fi
fi

echo "Server built successfully!"
echo ""

# Function to send a request and display response
send_request() {
    local description="$1"
    local request="$2"
    
    echo "=== $description ==="
    echo "Request:"
    echo "$request"
    echo ""
    echo "Response:"
    echo "$request" | ./mcp-server 2>/dev/null
    echo ""
    echo "---"
    echo ""
}

# Test 1: Multiline JSON-RPC request
echo "Test 1: Multiline JSON-RPC Request"
send_request "Multiline Initialize" '{
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
      "name": "multiline-test",
      "version": "1.0.0"
    }
  }
}'

# Test 2: Batch JSON-RPC requests
echo "Test 2: Batch JSON-RPC Requests"
send_request "Batch Requests (Initialize + List Tools + Echo)" '[
  {
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
        "name": "batch-test",
        "version": "1.0.0"
      }
    }
  },
  {
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
  },
  {
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "echo",
      "arguments": {
        "text": "Hello from batch request!"
      }
    }
  }
]'

# Test 3: Using the example files
echo "Test 3: Using Multiline Example File"
send_request "Multiline File Test" "$(cat ../examples/multiline_test.json)"

echo "Test 4: Using Batch Example File"
send_request "Batch File Test" "$(cat ../examples/batch_test.json)"

echo "Multiline and Batch Testing completed!"
echo "=====================================" 