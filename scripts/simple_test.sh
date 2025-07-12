#!/bin/bash

echo "Simple MCP Server Test"
echo "======================"

# Build the server
echo "Building server..."
go build -o mcp-server ../cmd/server/main.go

if [ $? -ne 0 ]; then
    echo "Failed to build server"
    exit 1
fi

echo "Server built successfully!"
echo ""

# Test function
test_request() {
    local test_name="$1"
    local json_file="$2"
    
    echo "=== $test_name ==="
    echo "Sending request from: $json_file"
    echo "Response:"
    cat "$json_file" | ./mcp-server
    echo ""
    echo "---"
    echo ""
}

# Run tests
test_request "1. Initialize" "../examples/01_initialize.json"
test_request "2. List Resources" "../examples/02_list_resources.json"
test_request "3. List Tools" "../examples/03_list_tools.json"
test_request "4. Call Echo Tool" "../examples/04_echo_tool.json"

echo "Test completed!" 