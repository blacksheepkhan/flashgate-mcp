package tools

import (
	"context"
	"encoding/json"

	"github.com/blacksheepkhan/flashgate-mcp/internal/fs"
	"github.com/blacksheepkhan/flashgate-mcp/internal/protocol"
)

const mkdirToolName = "mkdir"

// MkdirTool exposes directory creation as an MCP tool.
type MkdirTool struct {
	filesystem fs.FileSystem
}

// NewMkdirTool creates a new mkdir tool.
func NewMkdirTool(filesystem fs.FileSystem) *MkdirTool {
	return &MkdirTool{
		filesystem: filesystem,
	}
}

// Name returns the tool name.
func (t *MkdirTool) Name() string {
	return mkdirToolName
}

// Title returns the human-readable tool title.
func (t *MkdirTool) Title() string {
	return "Make Directory"
}

// Description returns the tool description.
func (t *MkdirTool) Description() string {
	return "Creates a directory below the configured filesystem root."
}

// InputSchema returns the JSON schema for this tool.
func (t *MkdirTool) InputSchema() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Relative directory path below the configured filesystem root.",
			},
		},
		"required":             []string{"path"},
		"additionalProperties": false,
	}
}

// Definition returns the MCP tool definition.
func (t *MkdirTool) Definition() protocol.Tool {
	return protocol.Tool{
		Name:        t.Name(),
		Title:       t.Title(),
		Description: t.Description(),
		InputSchema: t.InputSchema(),
	}
}

// Execute creates the requested directory.
func (t *MkdirTool) Execute(_ context.Context, rawArguments json.RawMessage) (any, *protocol.Error) {
	var arguments mkdirArguments
	if err := json.Unmarshal(rawArguments, &arguments); err != nil {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "invalid mkdir arguments",
		}
	}

	if arguments.Path == "" {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "missing path",
		}
	}

	if err := t.filesystem.Mkdir(arguments.Path); err != nil {
		return nil, mapFilesystemError(err)
	}

	return mkdirResult{
		Path:    arguments.Path,
		Created: true,
	}, nil
}

type mkdirArguments struct {
	Path string `json:"path"`
}

type mkdirResult struct {
	Path    string `json:"path"`
	Created bool   `json:"created"`
}
