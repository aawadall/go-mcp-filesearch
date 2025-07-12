package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aawadall/go-mcp-filesearch/internal/models"
)

// HTTPMCPClient represents a client for the HTTP-based MCP server
type HTTPMCPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPMCPClient creates a new HTTP MCP client
func NewHTTPMCPClient(baseURL string) *HTTPMCPClient {
	return &HTTPMCPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendRequest sends a JSON-RPC request to the MCP server
func (c *HTTPMCPClient) SendRequest(req models.JSONRPCRequest) (*models.JSONRPCResponse, error) {
	// Marshal request to JSON
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.baseURL+"/mcp", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
	}

	// Parse JSON-RPC response
	var mcpResp models.JSONRPCResponse
	if err := json.Unmarshal(body, &mcpResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &mcpResp, nil
}

// SendBatchRequest sends a batch of JSON-RPC requests
func (c *HTTPMCPClient) SendBatchRequest(requests []models.JSONRPCRequest) ([]models.JSONRPCResponse, error) {
	// Marshal batch request to JSON
	reqBytes, err := json.Marshal(requests)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal batch request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.baseURL+"/mcp", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
	}

	// Parse batch JSON-RPC response
	var responses []models.JSONRPCResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch response: %w", err)
	}

	return responses, nil
}

// Initialize initializes the MCP connection
func (c *HTTPMCPClient) Initialize(clientName, clientVersion string) error {
	req := models.JSONRPCRequest{
		JSONRPC: models.JSONRPCVersion,
		ID:      1,
		Method:  "initialize",
		Params: models.InitializeParams{
			ProtocolVersion: models.MCPProtocolVersion,
			Capabilities:    map[string]interface{}{},
			ClientInfo: models.ClientInfo{
				Name:    clientName,
				Version: clientVersion,
			},
		},
	}

	resp, err := c.SendRequest(req)
	if err != nil {
		return fmt.Errorf("initialize failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("initialize error: %s", resp.Error.Message)
	}

	fmt.Printf("Initialized MCP connection: %+v\n", resp.Result)
	return nil
}

// ListTools lists available tools
func (c *HTTPMCPClient) ListTools() error {
	req := models.JSONRPCRequest{
		JSONRPC: models.JSONRPCVersion,
		ID:      2,
		Method:  "tools/list",
		Params:  map[string]interface{}{},
	}

	resp, err := c.SendRequest(req)
	if err != nil {
		return fmt.Errorf("list tools failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("list tools error: %s", resp.Error.Message)
	}

	fmt.Printf("Available tools: %+v\n", resp.Result)
	return nil
}

// CallTool calls a specific tool
func (c *HTTPMCPClient) CallTool(name string, arguments map[string]interface{}) error {
	req := models.JSONRPCRequest{
		JSONRPC: models.JSONRPCVersion,
		ID:      3,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":      name,
			"arguments": arguments,
		},
	}

	resp, err := c.SendRequest(req)
	if err != nil {
		return fmt.Errorf("call tool failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("call tool error: %s", resp.Error.Message)
	}

	fmt.Printf("Tool result: %+v\n", resp.Result)
	return nil
}

func main() {
	// Create client
	client := NewHTTPMCPClient("http://localhost:8080")

	fmt.Println("=== HTTP MCP Client Example ===")

	// Initialize connection
	fmt.Println("\n1. Initializing MCP connection...")
	if err := client.Initialize("http-client-example", "1.0.0"); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// List tools
	fmt.Println("\n2. Listing available tools...")
	if err := client.ListTools(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Call echo tool
	fmt.Println("\n3. Calling echo tool...")
	if err := client.CallTool("echo", map[string]interface{}{
		"text": "Hello from Go HTTP client!",
	}); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Example batch request
	fmt.Println("\n4. Making batch request...")
	batchReq := []models.JSONRPCRequest{
		{
			JSONRPC: models.JSONRPCVersion,
			ID:      4,
			Method:  "tools/list",
			Params:  map[string]interface{}{},
		},
		{
			JSONRPC: models.JSONRPCVersion,
			ID:      5,
			Method:  "resources/list",
			Params:  map[string]interface{}{},
		},
	}

	responses, err := client.SendBatchRequest(batchReq)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Batch responses: %+v\n", responses)

	fmt.Println("\n=== Example completed ===")
}
