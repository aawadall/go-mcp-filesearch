#!/bin/bash

# Test script for the MCP server using batch requests
# This script sends all JSON-RPC requests as a single batch to maintain server state

echo "Starting MCP Server Test (Batch Version)"
echo "======================================="

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

# Create a batch request with all test commands
echo "=== Sending Batch Request ==="
echo "Request: Batch of all test commands"
echo "Response:"

# Send all requests as a batch
cat << 'EOF' | ./mcp-server 2>/dev/null
[
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
        "name": "test-client",
        "version": "1.0.0"
      }
    }
  },
  {
    "jsonrpc": "2.0",
    "id": 2,
    "method": "resources/list",
    "params": {}
  },
  {
    "jsonrpc": "2.0",
    "id": 3,
    "method": "resources/read",
    "params": {
      "uri": "example://test"
    }
  },
  {
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/list",
    "params": {}
  },
  {
    "jsonrpc": "2.0",
    "id": 5,
    "method": "tools/call",
    "params": {
      "name": "echo",
      "arguments": {
        "text": "Hello, World!"
      }
    }
  },
  {
    "jsonrpc": "2.0",
    "id": 6,
    "method": "tools/call",
    "params": {
      "name": "echo",
      "arguments": {
        "text": "Testing the MCP server!"
      }
    }
  },
  {
    "jsonrpc": "2.0",
    "id": 7,
    "method": "invalid/method",
    "params": {}
  },
  {
    "jsonrpc": "2.0",
    "id": 8,
    "method": "tools/call",
    "params": {
      "name": "nonexistent_tool",
      "arguments": {}
    }
  },
  {
    "jsonrpc": "2.0",
    "id": 9,
    "method": "tools/call",
    "params": {
      "name": "echo",
      "arguments": {}
    }
  }
]
EOF

echo ""
echo "Test completed!"
echo "===============" 