package tools

import (
	"context"
	"encoding/json"

	"github.com/blacksheepkhan/fileserver-mcp/internal/fs"
	"github.com/blacksheepkhan/fileserver-mcp/internal/protocol"
)

const movePathToolName = "move_path"

// MovePathTool exposes filesystem move operations as an MCP tool.
type MovePathTool struct {
	filesystem fs.FileSystem
}

// NewMovePathTool creates a new move_path tool.
func NewMovePathTool(filesystem fs.FileSystem) *MovePathTool {
	return &MovePathTool{
		filesystem: filesystem,
	}
}

// Name returns the tool name.
func (t *MovePathTool) Name() string {
	return movePathToolName
}

// Title returns the human-readable tool title.
func (t *MovePathTool) Title() string {
	return "Move Path"
}

// Description returns the tool description.
func (t *MovePathTool) Description() string {
	return "Moves a file or directory below the configured filesystem root."
}

// InputSchema returns the JSON schema for this tool.
func (t *MovePathTool) InputSchema() any {
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
func (t *MovePathTool) Definition() protocol.Tool {
	return protocol.Tool{
		Name:        t.Name(),
		Title:       t.Title(),
		Description: t.Description(),
		InputSchema: t.InputSchema(),
	}
}

// Execute moves the requested source path to the target path.
func (t *MovePathTool) Execute(_ context.Context, rawArguments json.RawMessage) (any, *protocol.Error) {
	var arguments movePathArguments
	if err := json.Unmarshal(rawArguments, &arguments); err != nil {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "invalid move_path arguments",
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

	if err := t.filesystem.Move(arguments.Source, arguments.Target, arguments.Overwrite); err != nil {
		return nil, mapFilesystemError(err)
	}

	return movePathResult{
		Source: arguments.Source,
		Target: arguments.Target,
		Moved:  true,
	}, nil
}

type movePathArguments struct {
	Source    string `json:"source"`
	Target    string `json:"target"`
	Overwrite bool   `json:"overwrite,omitempty"`
}

type movePathResult struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Moved  bool   `json:"moved"`
}
