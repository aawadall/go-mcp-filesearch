// Package server implements HTTP-based MCP server functionality.
// It provides HTTP transport for the Model Context Protocol (MCP) server.
package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aawadall/go-mcp-filesearch/internal/models"
)

// HTTPMCPServer wraps an MCPServer to provide HTTP transport
type HTTPMCPServer struct {
	mcpServer *MCPServer
	mux       *http.ServeMux
}

// NewHTTPMCPServer creates a new HTTP MCP server that wraps the given MCP server
func NewHTTPMCPServer(mcpServer *MCPServer) *HTTPMCPServer {
	httpServer := &HTTPMCPServer{
		mcpServer: mcpServer,
		mux:       http.NewServeMux(),
	}

	// Set up routes
	httpServer.setupRoutes()

	return httpServer
}

// setupRoutes configures all the HTTP routes for the MCP server
func (h *HTTPMCPServer) setupRoutes() {
	// Main MCP endpoint - handles all MCP protocol requests
	h.mux.HandleFunc("/mcp", h.handleMCPRequest)

	// Health check endpoint
	h.mux.HandleFunc("/health", h.HealthCheckHandler)

	// Server information endpoint
	h.mux.HandleFunc("/info", h.InfoHandler)

	// Root endpoint with basic info
	h.mux.HandleFunc("/", h.RootHandler)
}

// ServeHTTP delegates to the underlying mux
func (h *HTTPMCPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for web clients
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Delegate to mux
	h.mux.ServeHTTP(w, r)
}

// handleMCPRequest handles the main MCP protocol requests
func (h *HTTPMCPServer) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests for MCP operations
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type for JSON responses
	w.Header().Set("Content-Type", "application/json")

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the request
	var response interface{}

	// Try to parse as single JSON-RPC request
	var req models.JSONRPCRequest
	if err := json.Unmarshal(body, &req); err == nil {
		// Single request
		mcpResponse := h.mcpServer.handleRequest(req)
		response = mcpResponse
	} else {
		// Try to parse as batch request
		var requests []models.JSONRPCRequest
		if err := json.Unmarshal(body, &requests); err == nil {
			// Batch request
			responses := h.mcpServer.handleBatchRequest(requests)
			response = responses
		} else {
			// Invalid JSON
			response = models.JSONRPCResponse{
				JSONRPC: models.JSONRPCVersion,
				ID:      nil,
				Error: &models.JSONRPCError{
					Code:    models.ErrCodeParseError,
					Message: "Invalid JSON-RPC request",
				},
			}
		}
	}

	// Encode and send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// RootHandler provides basic information about the server
func (h *HTTPMCPServer) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":        "MCP HTTP Server",
		"version":     models.ServerVersion,
		"protocol":    models.MCPProtocolVersion,
		"transport":   "http",
		"description": "HTTP-based Model Context Protocol server",
		"endpoints": map[string]string{
			"mcp":    "/mcp - Main MCP protocol endpoint",
			"health": "/health - Health check",
			"info":   "/info - Server information",
		},
		"usage": "Send POST requests to /mcp with JSON-RPC 2.0 formatted MCP requests",
	})
}

// HealthCheckHandler provides a simple health check endpoint
func (h *HTTPMCPServer) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"server": "mcp-http-server",
	})
}

// InfoHandler provides server information endpoint
func (h *HTTPMCPServer) InfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":        models.ServerName,
		"version":     models.ServerVersion,
		"protocol":    models.MCPProtocolVersion,
		"transport":   "http",
		"description": "HTTP-based MCP server",
	})
}
