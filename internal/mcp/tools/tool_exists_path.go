package tools

import (
	"context"
	"encoding/json"

	"github.com/blacksheepkhan/fileserver-mcp/internal/fs"
	"github.com/blacksheepkhan/fileserver-mcp/internal/protocol"
)

const existsPathToolName = "exists_path"

// ExistsPathTool checks whether a path exists below the configured filesystem root.
type ExistsPathTool struct {
	filesystem fs.FileSystem
}

// NewExistsPathTool creates a new exists_path tool.
func NewExistsPathTool(filesystem fs.FileSystem) *ExistsPathTool {
	return &ExistsPathTool{
		filesystem: filesystem,
	}
}

// Name returns the tool name.
func (t *ExistsPathTool) Name() string {
	return existsPathToolName
}

// Title returns the human-readable tool title.
func (t *ExistsPathTool) Title() string {
	return "Exists Path"
}

// Description returns the tool description.
func (t *ExistsPathTool) Description() string {
	return "Checks whether a file or directory exists below the configured filesystem root."
}

// InputSchema returns the JSON schema for this tool.
func (t *ExistsPathTool) InputSchema() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Relative file or directory path below the configured filesystem root.",
			},
		},
		"required":             []string{"path"},
		"additionalProperties": false,
	}
}

// Definition returns the MCP tool definition.
func (t *ExistsPathTool) Definition() protocol.Tool {
	return protocol.Tool{
		Name:        t.Name(),
		Title:       t.Title(),
		Description: t.Description(),
		InputSchema: t.InputSchema(),
	}
}

// Execute checks whether the requested path exists.
func (t *ExistsPathTool) Execute(_ context.Context, rawArguments json.RawMessage) (any, *protocol.Error) {
	var arguments existsPathArguments
	if err := json.Unmarshal(rawArguments, &arguments); err != nil {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "invalid exists_path arguments",
		}
	}

	if arguments.Path == "" {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "missing path",
		}
	}

	exists, err := t.filesystem.Exists(arguments.Path)
	if err != nil {
		return nil, mapFilesystemError(err)
	}

	return existsPathResult{
		Exists: exists,
	}, nil
}

type existsPathArguments struct {
	Path string `json:"path"`
}

type existsPathResult struct {
	Exists bool `json:"exists"`
}
