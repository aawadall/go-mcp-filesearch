// Package server implements the Model Context Protocol (MCP) server functionality.
// It provides the core server logic for handling MCP requests, managing resources and tools,
// and communicating with MCP clients via JSON-RPC 2.0 over stdin/stdout.
package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aawadall/go-mcp-filesearch/internal/models"
)

// MCPServer represents an MCP server instance that handles client requests
// and manages resources and tools.
type MCPServer struct {
	initialized bool
	resources   []models.Resource
	tools       []models.Tool
}

// NewMCPServer creates and returns a new MCPServer instance with default
// resources and tools configured.
func NewMCPServer() *MCPServer {
	server := &MCPServer{
		resources: []models.Resource{
			{
				URI:         "example://test",
				Name:        "Test Resource",
				Description: "A simple test resource",
				MimeType:    "text/plain",
			},
		},
		tools: []models.Tool{
			{
				Name:        "echo",
				Description: "Echo back the input text",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"text": map[string]interface{}{
							"type":        "string",
							"description": "Text to echo back",
						},
					},
					"required": []string{"text"},
				},
			},
		},
	}
	return server
}

// handleInitialize processes the initialize method request and returns
// server capabilities and information.
func (s *MCPServer) handleInitialize(params interface{}) (interface{}, error) {
	s.initialized = true

	return models.InitializeResult{
		ProtocolVersion: models.MCPProtocolVersion,
		Capabilities: models.ServerCapabilities{
			Resources: map[string]interface{}{
				"subscribe":   true,
				"listChanged": true,
			},
			Tools: map[string]interface{}{
				"listChanged": true,
			},
		},
		ServerInfo: models.ServerInfo{
			Name:    models.ServerName,
			Version: models.ServerVersion,
		},
	}, nil
}

// handleListResources returns the list of available resources.
func (s *MCPServer) handleListResources(params interface{}) (interface{}, error) {
	if !s.initialized {
		return nil, fmt.Errorf("server not initialized")
	}

	return map[string]interface{}{
		"resources": s.resources,
	}, nil
}

// handleListTools returns the list of available tools.
func (s *MCPServer) handleListTools(params interface{}) (interface{}, error) {
	if !s.initialized {
		return nil, fmt.Errorf("server not initialized")
	}

	return map[string]interface{}{
		"tools": s.tools,
	}, nil
}

// handleReadResource reads and returns the contents of a specified resource.
func (s *MCPServer) handleReadResource(params interface{}) (interface{}, error) {
	if !s.initialized {
		return nil, fmt.Errorf("server not initialized")
	}

	return map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":      "example://test",
				"mimeType": "text/plain",
				"text":     "This is a test content from the MCP server",
			},
		},
	}, nil
}

// handleCallTool executes a specific tool with the provided arguments.
func (s *MCPServer) handleCallTool(params interface{}) (interface{}, error) {
	if !s.initialized {
		return nil, fmt.Errorf("server not initialized")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid params")
	}

	name, ok := paramsMap["name"].(string)
	if !ok {
		return nil, fmt.Errorf("tool name required")
	}

	switch name {
	case "echo":
		args, ok := paramsMap["arguments"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("arguments required")
		}

		text, ok := args["text"].(string)
		if !ok {
			return nil, fmt.Errorf("text argument required")
		}

		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Echo: %s", text),
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

// handleBatchRequest processes a batch of JSON-RPC requests and returns an array of responses.
func (s *MCPServer) handleBatchRequest(requests []models.JSONRPCRequest) []models.JSONRPCResponse {
	responses := make([]models.JSONRPCResponse, len(requests))

	for i, req := range requests {
		responses[i] = s.handleRequest(req)
	}

	return responses
}

// handleRequest routes incoming JSON-RPC requests to the appropriate handler method.
func (s *MCPServer) handleRequest(req models.JSONRPCRequest) models.JSONRPCResponse {
	var result interface{}
	var err error

	switch req.Method {
	case "initialize":
		result, err = s.handleInitialize(req.Params)
	case "resources/list":
		result, err = s.handleListResources(req.Params)
	case "resources/read":
		result, err = s.handleReadResource(req.Params)
	case "tools/list":
		result, err = s.handleListTools(req.Params)
	case "tools/call":
		result, err = s.handleCallTool(req.Params)
	default:
		err = fmt.Errorf("method not found: %s", req.Method)
	}

	response := models.JSONRPCResponse{
		JSONRPC: models.JSONRPCVersion,
		ID:      req.ID,
	}

	if err != nil {
		response.Error = &models.JSONRPCError{
			Code:    models.ErrCodeMethodNotFound, // method not found
			Message: err.Error(),
		}
	} else {
		response.Result = result
	}

	return response
}

// Run starts the MCP server and begins listening for JSON-RPC requests on stdin.
// The server processes requests and supports both single-line and multiline JSON-RPC messages.
func (s *MCPServer) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Add the line to our buffer
		buffer.WriteString(line)
		buffer.WriteString("\n")

		// Try to parse the accumulated content as JSON
		content := buffer.String()
		content = strings.TrimSpace(content)

		if content == "" {
			continue
		}

		// First, try to parse as a single JSON-RPC request
		var req models.JSONRPCRequest
		if err := json.Unmarshal([]byte(content), &req); err == nil {
			// Single request - process it
			response := s.handleRequest(req)
			if respBytes, err := json.Marshal(response); err == nil {
				fmt.Println(string(respBytes))
			}
			buffer.Reset()
			continue
		}

		// If single request parsing failed, try parsing as a batch request
		var requests []models.JSONRPCRequest
		if err := json.Unmarshal([]byte(content), &requests); err == nil {
			// Batch request - process all requests
			responses := s.handleBatchRequest(requests)
			for _, response := range responses {
				if respBytes, err := json.Marshal(response); err == nil {
					fmt.Println(string(respBytes))
				}
			}
			buffer.Reset()
			continue
		}

		// If both parsing attempts failed, continue accumulating lines
		// This allows for multiline JSON input
	}

	// Handle any remaining content in buffer (incomplete JSON)
	if remaining := strings.TrimSpace(buffer.String()); remaining != "" {
		errResp := models.JSONRPCResponse{
			JSONRPC: models.JSONRPCVersion,
			ID:      nil,
			Error: &models.JSONRPCError{
				Code:    models.ErrCodeParseError,
				Message: "Incomplete JSON-RPC request",
			},
		}
		if respBytes, err := json.Marshal(errResp); err == nil {
			fmt.Println(string(respBytes))
		}
	}
}
