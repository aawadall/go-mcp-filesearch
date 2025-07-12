package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

)

// JSON-RPC 2.0 structures
type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID interface{} `json:"id"`
	Method string `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID interface{} `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error *JSONRPCError `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

// MCP specific structures
type InitializeParams struct {
	ProtocolVersion string `json:"protocolVersion"`
	Capabilities map[string]interface{} `json:"capabilities"`
	ClientInfo ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name string `json:"name"`
	Version string `json:"version"`
}

type ServerInfo struct {
	Name string `json:"name"`
	Version string `json:"version"`	
}

type ServerCapabilities struct {
	Resources map[string]interface{} `json:"resources,omitempty"`
	Tools map[string]interface{} `json:"tools,omitempty"`
	Prompts map[string]interface{} `json:"prompts,omitempty"`
}

type InitializeResult struct {
	ProtocolVersion string `json:"protocolVersion"`
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo ServerInfo `json:"serverInfo"`
}

type Resource struct {
	URI string `json:"uri"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

type Tool struct {
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type MCPServer struct {
	initialized bool
	resources []Resource
	tools []Tool
}

func NewMCPServer() *MCPServer {
	server := &MCPServer{
		resources: []Resource{
			{
				URI: "example://test",
				Name: "Test Resource",
				Description: "A simple test resource",
				MimeType: "text/plain",
			},
		},
		tools: []Tool {
			{
				Name: "echo",
				Description: "Echo back the input text",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"text": map[string]interface{}{
							"type": "string",
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

func (s *MCPServer) handleInitialize(params interface{}) (interface{}, error) {
	s.initialized = true

	return InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Resources: map[string]interface{}{
				"subscribe": true,
				"listChanged": true,
			},
			Tools: map[string]interface{}{
				"listChanged": true,
			},
		},
		ServerInfo: ServerInfo{
			Name: "simple-mcp-server",
			Version: "1.0.0",
		},
	},nil
}

func (s *MCPServer) handleListResources(params interface{}) (interface{}, error) {
	if !s.initialized {
		return nil, fmt.Errorf("server not initialized")
	}

	return map[string]interface{}{
		"resources": s.resources,
	}, nil
}

func (s *MCPServer) handleListTools(params interface{}) (interface{}, error) {
	if !s.initialized {
		return nil, fmt.Errorf("server not initialized")
	}

	return map[string]interface{}{
		"tools": s.tools,
	}, nil
}

func (s *MCPServer) handleReadResource(params interface{}) (interface{}, error) {
	if !s.initialized {
		return nil, fmt.Errorf("server not initialized")
	}

	return map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri": "example://test",
				"mimeType": "text/plain",
				"text": "This is a test content from the MCP server",
			},
		},
	},nil
}

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

func (s *MCPServer) handleRequest(req JSONRPCRequest) JSONRPCResponse {
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

	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID: req.ID,
	}

	if err != nil {
		response.Error = &JSONRPCError{
			Code: -32601, // method not found
			Message: err.Error(),
		}
	} else {
		response.Result = result
	}

	return response
}

func (s *MCPServer) Run() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			errResp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID: nil,
				Error: &JSONRPCError{
					Code: -32700, // Parse error
					Message: "Parse error",
				},

			}
			if respBytes, err := json.Marshal(errResp); err == nil {
				fmt.Println(string(respBytes))
			}
			continue
		}

		response := s.handleRequest(req)

		if respBytes, err := json.Marshal(response); err == nil {
			fmt.Println(string(respBytes))
		}
	}
}

func main() {
	server := NewMCPServer()
	server.Run()
}
