# Go MCP File Search Server

A Model Context Protocol (MCP) server implementation in Go that provides file search capabilities and basic tool functionality.

## Overview

This project implements an MCP server that follows the [Model Context Protocol specification](https://modelcontextprotocol.io/). The server provides:

- **Resources**: File-based resources that can be listed and read
- **Tools**: Executable tools that can be called by MCP clients
- **JSON-RPC 2.0 Communication**: Standard protocol communication over stdin/stdout

## Features

### Resources
- **Test Resource**: A simple example resource with URI `example://test`
- **Resource Listing**: Lists available resources via `resources/list`
- **Resource Reading**: Reads resource contents via `resources/read`

### Tools
- **Echo Tool**: A simple tool that echoes back input text
  - Input: `{"text": "string"}`
  - Output: `{"content": [{"type": "text", "text": "Echo: <input>"}]}`

### MCP Protocol Support
- **Initialize**: Handles client initialization with protocol version `2024-11-05`
- **Capabilities**: Supports resource subscription, list changes, and tool list changes
- **Error Handling**: Proper JSON-RPC error responses with appropriate error codes

## Project Structure

```
go-mcp-filesearch/
├── cmd/
│   └── server/
│       └── main.go          # Main entry point
├── internal/
│   ├── filesearch/          # File search functionality (placeholder)
│   │   ├── handler.go
│   │   └── registry.go
│   └── server/
│       └── server.go        # MCP server implementation
├── go.mod                   # Go module definition
└── README.md               # This file
```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/aawadall/go-mcp-filesearch.git
cd go-mcp-filesearch
```

2. Build the server:
```bash
go build -o mcp-server cmd/server/main.go
```

## Usage

### Running the Server

The server communicates via stdin/stdout using JSON-RPC 2.0:

```bash
./mcp-server
```

### Example MCP Client Integration

The server can be integrated with MCP clients like Claude Desktop or other MCP-compatible applications. The server expects JSON-RPC messages on stdin and responds with JSON-RPC messages on stdout.

### Example Request/Response

**Initialize Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "test-client",
      "version": "1.0.0"
    }
  }
}
```

**Initialize Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
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
    "serverInfo": {
      "name": "simple-mcp-server",
      "version": "1.0.0"
    }
  }
}
```

**Tool Call Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "echo",
    "arguments": {
      "text": "Hello, World!"
    }
  }
}
```

**Tool Call Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Echo: Hello, World!"
      }
    ]
  }
}
```

## Supported Methods

- `initialize` - Initialize the MCP connection
- `resources/list` - List available resources
- `resources/read` - Read resource contents
- `tools/list` - List available tools
- `tools/call` - Call a specific tool

## Development

### Adding New Tools

To add a new tool, modify the `NewMCPServer()` function in `internal/server/server.go`:

```go
tools: []Tool{
    {
        Name:        "your-tool-name",
        Description: "Description of your tool",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "param1": map[string]interface{}{
                    "type":        "string",
                    "description": "Parameter description",
                },
            },
            "required": []string{"param1"},
        },
    },
}
```

Then add the tool handling logic in the `handleCallTool` method.

### Adding New Resources

To add new resources, modify the `resources` slice in `NewMCPServer()`:

```go
resources: []Resource{
    {
        URI:         "your://resource-uri",
        Name:        "Your Resource Name",
        Description: "Description of your resource",
        MimeType:    "text/plain",
    },
}
```

## Requirements

- Go 1.23.0 or later
- No external dependencies (uses only standard library)

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]
