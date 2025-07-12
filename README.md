# Go MCP File Search Server

A Model Context Protocol (MCP) server implementation in Go that provides file search capabilities and basic tool functionality.

## Overview

This project implements an MCP server that follows the [Model Context Protocol specification](https://modelcontextprotocol.io/). The server provides:

- **Resources**: File-based resources that can be listed and read
- **Tools**: Executable tools that can be called by MCP clients
- **Multiple Transport Options**: 
  - **stdin/stdout**: Standard protocol communication over stdin/stdout
  - **HTTP**: RESTful HTTP API for web-based clients

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
- **Multiline Support**: Can parse JSON-RPC requests spanning multiple lines
- **Batch Processing**: Supports processing multiple JSON-RPC requests in a single input

## Project Structure

```
go-mcp-filesearch/
├── cmd/
│   ├── server/
│   │   └── main.go          # stdin/stdout MCP server entry point
│   └── http-server/
│       └── main.go          # HTTP MCP server entry point
├── internal/
│   ├── filesearch/          # File search functionality (placeholder)
│   │   ├── handler.go
│   │   └── registry.go
│   ├── models/
│   │   └── mcp.go          # MCP and JSON-RPC data structures and constants
│   └── server/
│       ├── server.go        # MCP server implementation and business logic
│       └── http_server.go   # HTTP transport layer for MCP server
├── examples/
│   └── http_client.go       # Example HTTP client implementation
├── scripts/
│   ├── test_http_server.sh  # HTTP server test script
│   └── ...                  # Other test scripts
├── go.mod                   # Go module definition
└── README.md               # This file
```

### Code Organization

The project follows Go best practices with clear separation of concerns:

- **`internal/models/`**: Contains all data structures, constants, and type definitions
  - JSON-RPC 2.0 structures (`JSONRPCRequest`, `JSONRPCResponse`, `JSONRPCError`)
  - MCP-specific structures (`InitializeParams`, `ServerInfo`, `Resource`, `Tool`)
  - Protocol constants and error codes
- **`internal/server/`**: Contains the server implementation and business logic
  - `MCPServer` struct and its methods
  - Request handling and routing logic
  - Tool and resource management
- **`cmd/server/`**: Contains the main application entry point

## Installation

1. Clone the repository:
```bash
git clone https://github.com/aawadall/go-mcp-filesearch.git
cd go-mcp-filesearch
```

2. Build the servers:
```bash
# Build stdin/stdout MCP server
go build -o mcp-server cmd/server/main.go

# Build HTTP MCP server
go build -o mcp-http-server cmd/http-server/main.go
```

## Usage

### Running the Servers

#### stdin/stdout MCP Server

The traditional MCP server communicates via stdin/stdout using JSON-RPC 2.0:

```bash
./mcp-server
```

#### HTTP MCP Server

The HTTP-based MCP server provides a RESTful API for web clients:

```bash
# Run on default port 8080
./mcp-http-server

# Run on custom port
MCP_HTTP_PORT=9000 ./mcp-http-server
```

The HTTP server provides the following endpoints:

- `GET /` - Server information and usage guide
- `GET /health` - Health check endpoint
- `GET /info` - Server details and capabilities
- `POST /mcp` - Main MCP protocol endpoint (accepts JSON-RPC 2.0 requests)

#### Testing the HTTP Server

Use the provided test script to verify the HTTP server functionality:

```bash
# Start the HTTP server in one terminal
./mcp-http-server

# In another terminal, run the test script
./scripts/test_http_server.sh
```

#### Using the Go HTTP Client

The project includes an example Go client for the HTTP server:

```bash
# Build and run the example client
go run examples/http_client.go
```

### Input Formats

The server supports three input formats:

1. **Single-line JSON-RPC**: Compact format on one line
2. **Multiline JSON-RPC**: Pretty-printed JSON spanning multiple lines
3. **Batch JSON-RPC**: Array of multiple JSON-RPC requests

### Testing Multiline and Batch Support

Use the provided test script to verify multiline and batch parsing:

```bash
./scripts/test_multiline.sh
```

### Example MCP Client Integration

#### stdin/stdout Transport

The traditional server can be integrated with MCP clients like Claude Desktop or other MCP-compatible applications. The server expects JSON-RPC messages on stdin and responds with JSON-RPC messages on stdout.

#### HTTP Transport

The HTTP server can be integrated with web-based clients, mobile apps, or any HTTP client. Send POST requests to `/mcp` with JSON-RPC 2.0 formatted requests:

```bash
# Example using curl
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {
        "name": "curl-client",
        "version": "1.0.0"
      }
    }
  }'
```

### Example Request/Response

**Single-Line Initialize Request:**
```json
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}
```

**Multiline Initialize Request:**
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

**Batch Request (Multiple Commands):**
```json
[
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
  },
  {
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
  }
]
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

### Code Structure Principles

The codebase follows these organizational principles:

1. **Separation of Concerns**: Data models are separate from business logic
2. **Single Responsibility**: Each package has a clear, focused purpose
3. **Dependency Management**: Internal packages use explicit imports
4. **Constants Management**: All magic numbers and strings are defined as constants

### Adding New Tools

To add a new tool, modify the `NewMCPServer()` function in `internal/server/server.go`:

```go
tools: []models.Tool{
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
resources: []models.Resource{
    {
        URI:         "your://resource-uri",
        Name:        "Your Resource Name",
        Description: "Description of your resource",
        MimeType:    "text/plain",
    },
}
```

### Adding New Data Structures

When adding new MCP or JSON-RPC structures, add them to `internal/models/mcp.go`:

```go
type NewStructure struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2,omitempty"`
}
```

### Error Handling

The server uses standardized JSON-RPC error codes defined in `internal/models/mcp.go`:

- `ErrCodeMethodNotFound` (-32601): Method not found
- `ErrCodeParseError` (-32700): Parse error

## Requirements

- Go 1.23.0 or later
- No external dependencies (uses only standard library)

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]
