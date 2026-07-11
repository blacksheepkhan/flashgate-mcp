package tools

import (
	"context"
	"encoding/json"

	"github.com/blacksheepkhan/flashgate-mcp/internal/fs"
	"github.com/blacksheepkhan/flashgate-mcp/internal/protocol"
)

const renamePathToolName = "rename_path"

// RenamePathTool exposes filesystem rename operations as an MCP tool.
type RenamePathTool struct {
	filesystem fs.FileSystem
}

// NewRenamePathTool creates a new rename_path tool.
func NewRenamePathTool(filesystem fs.FileSystem) *RenamePathTool {
	return &RenamePathTool{
		filesystem: filesystem,
	}
}

// Name returns the tool name.
func (t *RenamePathTool) Name() string {
	return renamePathToolName
}

// Title returns the human-readable tool title.
func (t *RenamePathTool) Title() string {
	return "Rename Path"
}

// Description returns the tool description.
func (t *RenamePathTool) Description() string {
	return "Renames a file or directory below the configured filesystem root."
}

// InputSchema returns the JSON schema for this tool.
func (t *RenamePathTool) InputSchema() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"source": map[string]any{
				"type":        "string",
				"description": "Relative source file or directory path below the configured filesystem root.",
			},
			"target": map[string]any{
				"type":        "string",
				"description": "Relative target file or directory path below the configured filesystem root.",
			},
			"overwrite": map[string]any{
				"type":        "boolean",
				"description": "Whether an existing target may be overwritten. Defaults to false.",
			},
		},
		"required":             []string{"source", "target"},
		"additionalProperties": false,
	}
}

// Definition returns the MCP tool definition.
func (t *RenamePathTool) Definition() protocol.Tool {
	return protocol.Tool{
		Name:        t.Name(),
		Title:       t.Title(),
		Description: t.Description(),
		InputSchema: t.InputSchema(),
	}
}

// Execute renames the requested source path to the target path.
func (t *RenamePathTool) Execute(_ context.Context, rawArguments json.RawMessage) (any, *protocol.Error) {
	var arguments renamePathArguments
	if err := json.Unmarshal(rawArguments, &arguments); err != nil {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "invalid rename_path arguments",
		}
	}

	if arguments.Source == "" {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "missing source",
		}
	}

	if arguments.Target == "" {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "missing target",
		}
	}

	if err := t.filesystem.Rename(arguments.Source, arguments.Target, arguments.Overwrite); err != nil {
		return nil, mapFilesystemError(err)
	}

	return renamePathResult{
		Source:  arguments.Source,
		Target:  arguments.Target,
		Renamed: true,
	}, nil
}

type renamePathArguments struct {
	Source    string `json:"source"`
	Target    string `json:"target"`
	Overwrite bool   `json:"overwrite,omitempty"`
}

type renamePathResult struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	Renamed bool   `json:"renamed"`
}
