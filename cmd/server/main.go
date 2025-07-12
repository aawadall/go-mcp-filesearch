package main

import (
	"github.com/aawadall/go-mcp-filesearch/internal/server"
)

func main() {
	mcpServer := server.NewMCPServer()
	mcpServer.Run()
}
