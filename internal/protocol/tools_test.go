package protocol

import (
	"encoding/json"
	"testing"
)

func TestToolMarshal(t *testing.T) {
	t.Parallel()

	tool := Tool{
		Name:        "list_files",
		Title:       "List Files",
		Description: "Lists files and directories.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type": "string",
				},
			},
			"required": []string{"path"},
		},
	}

	encoded, err := json.Marshal(tool)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}

	if decoded["name"] != "list_files" {
		t.Fatalf("expected name %q, got %v", "list_files", decoded["name"])
	}

	if decoded["title"] != "List Files" {
		t.Fatalf("expected title %q, got %v", "List Files", decoded["title"])
	}

	if decoded["description"] != "Lists files and directories." {
		t.Fatalf("unexpected description: %v", decoded["description"])
	}

	if _, ok := decoded["inputSchema"].(map[string]any); !ok {
		t.Fatalf("expected inputSchema object, got %#v", decoded["inputSchema"])
	}
}

func TestToolMarshalOmitsEmptyTitle(t *testing.T) {
	t.Parallel()

	tool := Tool{
		Name:        "list_files",
		Description: "Lists files and directories.",
		InputSchema: map[string]any{
			"type": "object",
		},
	}

	encoded, err := json.Marshal(tool)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}

	if _, exists := decoded["title"]; exists {
		t.Fatal("did not expect title field")
	}
}
