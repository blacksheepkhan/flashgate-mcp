package main

import (
	"github.com/blacksheepkhan/fileserver-mcp/internal/fs"
	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/tools"
)

type toolCapabilities struct {
	filesystemWrite bool
}

func capabilitiesFromReadOnly(readOnly bool) toolCapabilities {
	return toolCapabilities{
		filesystemWrite: !readOnly,
	}
}

func createToolRegistry(filesystem fs.FileSystem, maxFileSize int64, capabilities toolCapabilities) *tools.Registry {
	toolRegistry := tools.NewRegistry()
	toolRegistry.Register(tools.NewListFilesTool(filesystem))
	toolRegistry.Register(tools.NewReadFileTool(filesystem, maxFileSize))
	toolRegistry.Register(tools.NewStatPathTool(filesystem))
	toolRegistry.Register(tools.NewExistsPathTool(filesystem))

	if !capabilities.filesystemWrite {
		return toolRegistry
	}

	toolRegistry.Register(tools.NewWriteFileTool(filesystem))
	toolRegistry.Register(tools.NewMkdirTool(filesystem))
	toolRegistry.Register(tools.NewDeletePathTool(filesystem))
	toolRegistry.Register(tools.NewMovePathTool(filesystem))
	toolRegistry.Register(tools.NewCopyPathTool(filesystem))
	toolRegistry.Register(tools.NewRenamePathTool(filesystem))

	return toolRegistry
}
