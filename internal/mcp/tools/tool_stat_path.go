package tools

import (
	"context"
	"encoding/json"

	"github.com/blacksheepkhan/flashgate-mcp/internal/fs"
	"github.com/blacksheepkhan/flashgate-mcp/internal/protocol"
)

const statPathToolName = "stat_path"

// StatPathTool exposes filesystem metadata lookup as an MCP tool.
type StatPathTool struct {
	filesystem fs.FileSystem
}

// NewStatPathTool creates a new stat_path tool.
func NewStatPathTool(filesystem fs.FileSystem) *StatPathTool {
	return &StatPathTool{
		filesystem: filesystem,
	}
}

// Name returns the tool name.
func (t *StatPathTool) Name() string {
	return statPathToolName
}

// Title returns the human-readable tool title.
func (t *StatPathTool) Title() string {
	return "Stat Path"
}

// Description returns the tool description.
func (t *StatPathTool) Description() string {
	return "Returns metadata for a file or directory below the configured filesystem root."
}

// InputSchema returns the JSON schema for this tool.
func (t *StatPathTool) InputSchema() any {
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
func (t *StatPathTool) Definition() protocol.Tool {
	return protocol.Tool{
		Name:        t.Name(),
		Title:       t.Title(),
		Description: t.Description(),
		InputSchema: t.InputSchema(),
	}
}

// Execute returns metadata for the requested path.
func (t *StatPathTool) Execute(_ context.Context, rawArguments json.RawMessage) (any, *protocol.Error) {
	var arguments statPathArguments
	if err := json.Unmarshal(rawArguments, &arguments); err != nil {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "invalid stat_path arguments",
		}
	}

	if arguments.Path == "" {
		return nil, &protocol.Error{
			Code:    protocol.ErrInvalidParams,
			Message: "missing path",
		}
	}

	metadata, err := t.filesystem.Stat(arguments.Path)
	if err != nil {
		return nil, mapFilesystemError(err)
	}

	return statPathResult{
		Name:  metadata.Name,
		IsDir: metadata.IsDir,
		Size:  metadata.Size,
	}, nil
}

type statPathArguments struct {
	Path string `json:"path"`
}

type statPathResult struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}
