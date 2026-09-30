package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// RunMCPServer runs the MCP JSON-RPC 2.0 stdio server for deckPath
func RunMCPServer(deckPath string) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		line = []byte(strings.TrimSpace(string(line)))
		if len(line) == 0 {
			continue
		}

		var req jsonRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			continue
		}

		resp := handleMCPRequest(deckPath, req)
		if resp != nil {
			data, _ := json.Marshal(resp)
			fmt.Printf("%s\n", string(data))
		}
	}
}

func handleMCPRequest(deckPath string, req jsonRPCRequest) *jsonRPCResponse {
	switch req.Method {
	case "initialize":
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"serverInfo": map[string]string{
					"name":    "deckforge-mcp",
					"version": "1.2.0",
				},
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
			},
		}

	case "notifications/initialized":
		return nil

	case "tools/list":
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": GetMCPTools(),
			},
		}

	case "tools/call":
		var callParams struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		_ = json.Unmarshal(req.Params, &callParams)

		toolResult := ExecuteMCPTool(deckPath, callParams.Name, callParams.Arguments)
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"content": []map[string]string{
					{
						"type": "text",
						"text": toolResult,
					},
				},
			},
		}

	default:
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
			},
		}
	}
}

// ExecuteMCPTool dispatches an MCP tool call by name to its corresponding domain handler
func ExecuteMCPTool(deckPath, name string, args map[string]interface{}) string {
	switch name {
	case "deckforge_get_settings", "deckforge_update_settings":
		return handleSettingsTools(deckPath, name, args)

	case "deckforge_create_deck", "deckforge_get_deck_info", "deckforge_set_deck_meta":
		return handleDeckTools(deckPath, name, args)

	case "deckforge_list_slides", "deckforge_get_slide", "deckforge_update_slide", "deckforge_reorder_slides":
		return handleSlideTools(deckPath, name, args)

	case "deckforge_list_themes", "deckforge_set_theme", "deckforge_create_theme":
		return handleThemeTools(deckPath, name, args)

	case "deckforge_get_catalog", "deckforge_list_components", "deckforge_create_component", "deckforge_insert_component":
		return handleComponentTools(deckPath, name, args)

	case "deckforge_audit_tokens", "deckforge_apply_token_fix", "deckforge_get_styling_guide":
		return handleTokenTools(deckPath, name, args)

	default:
		return fmt.Sprintf("Unknown tool %s", name)
	}
}
