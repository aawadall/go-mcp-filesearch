// Package models provides data structures and constants for the Model Context Protocol (MCP) server.
// It contains all JSON-RPC 2.0 structures, MCP-specific types, and protocol constants used
// throughout the application.
package models

// Constants for JSON-RPC 2.0 protocol version
const (
	JSONRPCVersion = "2.0"
)

// Constants for MCP Protocol configuration
const (
	MCPProtocolVersion = "2024-11-05"
	ServerVersion      = "1.0.0"
	ServerName         = "simple-mcp-server"
)

// JSON-RPC 2.0 standard error codes
const (
	ErrCodeMethodNotFound = -32601 // Method not found
	ErrCodeParseError     = -32700 // Parse error
)

// JSON-RPC 2.0 structures for request/response communication

// JSONRPCRequest represents a JSON-RPC 2.0 request message
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response message
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error object
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP specific structures for protocol communication

// InitializeParams contains the parameters for the initialize method
type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
}

// ClientInfo contains information about the MCP client
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ServerInfo contains information about the MCP server
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ServerCapabilities describes the server's capabilities for resources, tools, and prompts
type ServerCapabilities struct {
	Resources map[string]interface{} `json:"resources,omitempty"`
	Tools     map[string]interface{} `json:"tools,omitempty"`
	Prompts   map[string]interface{} `json:"prompts,omitempty"`
}

// InitializeResult contains the result of the initialize method
type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

// Resource represents an MCP resource that can be listed and read
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// Tool represents an MCP tool that can be called by clients
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}
