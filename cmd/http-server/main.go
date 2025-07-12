package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aawadall/go-mcp-filesearch/internal/server"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("MCP_HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	// Create MCP server instance
	mcpServer := server.NewMCPServer()

	// Create HTTP server
	httpServer := server.NewHTTPMCPServer(mcpServer)

	// Start HTTP server
	log.Printf("Starting HTTP MCP server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, httpServer))
}
