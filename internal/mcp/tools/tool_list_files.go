package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blacksheepkhan/fileserver-mcp/internal/fs"
	"github.com/blacksheepkhan/fileserver-mcp/internal/protocol"
	"github.com/blacksheepkhan/fileserver-mcp/internal/security"
)

const listFilesToolName = "list_files"

// ListFilesTool exposes directory listing as an MCP tool.
type ListFilesTool struct {
	filesystem fs.FileSystem
}

// NewListFilesTool creates a new list_files tool.
func NewListFilesTool(filesystem fs.FileSystem) *ListFilesTool {
	return &ListFilesTool{
		filesystem: filesystem,
	}
}

// Name returns the tool name.
func (t *ListFilesTool) Name() string {
	return listFilesToolName
}

// Title returns the human-readable tool title.
func (t *ListFilesTool) Title() string {
	return "List Files"
}

// Description returns the tool description.
func (t *ListFilesTool) Description() string {
	return "Lists files and directories below the configured filesystem root."
}

// InputSchema returns the JSON schema for this tool.
func (t *ListFilesTool) InputSchema() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Relative path below the configured filesystem root. Defaults to '.'.",
			},
		},
		"additionalProperties": false,
	}
}

// Definition returns the MCP tool definition.
func (t *ListFilesTool) Definition() protocol.Tool {
	return protocol.Tool{
		Name:        t.Name(),
		Title:       t.Title(),
		Description: t.Description(),
		InputSchema: t.InputSchema(),
	}
}

// Execute lists files for the requested relative path.
func (t *ListFilesTool) Execute(_ context.Context, rawArguments json.RawMessage) (any, *protocol.Error) {
	var arguments listFilesArguments
	if len(rawArguments) > 0 {
		if err := json.Unmarshal(rawArguments, &arguments); err != nil {
			return nil, &protocol.Error{
				Code:    protocol.ErrInvalidParams,
				Message: "invalid list_files arguments",
			}
		}
	}

	path := arguments.Path
	if path == "" {
		path = "."
	}

	entries, err := t.filesystem.List(path)
	if err != nil {
		return nil, mapFilesystemError(err)
	}

	return listFilesResult{
		Entries: entries,
	}, nil
}

type listFilesArguments struct {
	Path string `json:"path,omitempty"`
}

type listFilesResult struct {
	Entries []fs.Entry `json:"entries"`
}

func mapFilesystemError(err error) *protocol.Error {
	code := protocol.ErrInternalError
	message := fmt.Sprintf("filesystem error: %v", err)

	if errors.Is(err, security.ErrAbsolutePath) ||
		errors.Is(err, security.ErrPathTraversal) ||
		errors.Is(err, security.ErrOutsideRoot) {
		code = protocol.ErrInvalidParams
		message = "filesystem error: invalid path"
	} else if errors.Is(err, fs.ErrPathIsDirectory) ||
		errors.Is(err, fs.ErrPathIsNotDirectory) ||
		errors.Is(err, fs.ErrFileTooLarge) ||
		errors.Is(err, fs.ErrFileExists) ||
		errors.Is(err, fs.ErrDirectoryNotEmpty) ||
		errors.Is(err, fs.ErrCopyDirectoryUnsupported) {
		code = protocol.ErrInvalidParams
	}

	return &protocol.Error{
		Code:    code,
		Message: message,
	}
}
