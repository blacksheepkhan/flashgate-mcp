package tools

import (
	"encoding/json"
	"testing"

	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/handlers"
)

func TestListHandlerReturnsRegisteredTools(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	registry.Register(&testTool{
		name:        "list_files",
		title:       "List Files",
		description: "Lists files and directories.",
		inputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type": "string",
				},
			},
		},
	})
	registry.Register(&testTool{
		name:        "read_file",
		title:       "Read File",
		description: "Reads a file.",
		inputSchema: map[string]any{
			"type": "object",
		},
	})

	handler := NewListHandler(registry)

	result, rpcErr := handler.Handle(handlers.Context{}, nil)
	if rpcErr != nil {
		t.Fatalf("expected no error, got %v", rpcErr)
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("expected marshalable result, got %v", err)
	}

	var decoded struct {
		Tools []struct {
			Name        string         `json:"name"`
			Title       string         `json:"title"`
			Description string         `json:"description"`
			InputSchema map[string]any `json:"inputSchema"`
		} `json:"tools"`
	}

	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid result JSON, got %v", err)
	}

	if len(decoded.Tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(decoded.Tools))
	}

	names := map[string]bool{}
	for _, tool := range decoded.Tools {
		names[tool.Name] = true

		if tool.Description == "" {
			t.Fatalf("expected description for tool %q", tool.Name)
		}

		if tool.InputSchema == nil {
			t.Fatalf("expected input schema for tool %q", tool.Name)
		}
	}

	if !names["list_files"] {
		t.Fatal("expected list_files tool")
	}

	if !names["read_file"] {
		t.Fatal("expected read_file tool")
	}
}

func TestListHandlerMethod(t *testing.T) {
	t.Parallel()

	handler := NewListHandler(NewRegistry())

	if handler.Method() != "tools/list" {
		t.Fatalf("expected tools/list, got %q", handler.Method())
	}
}
