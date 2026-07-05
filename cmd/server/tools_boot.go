package main

import (
	"github.com/blacksheepkhan/fileserver-mcp/internal/fs"
	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/tools"
)

func createToolRegistry(filesystem fs.FileSystem, maxFileSize int64) *tools.Registry {
	toolRegistry := tools.NewRegistry()
	toolRegistry.Register(tools.NewListFilesTool(filesystem))
	toolRegistry.Register(tools.NewReadFileTool(filesystem, maxFileSize))
	toolRegistry.Register(tools.NewStatPathTool(filesystem))
	toolRegistry.Register(tools.NewExistsPathTool(filesystem))
	toolRegistry.Register(tools.NewWriteFileTool(filesystem))
	toolRegistry.Register(tools.NewMkdirTool(filesystem))
	toolRegistry.Register(tools.NewDeletePathTool(filesystem))
	toolRegistry.Register(tools.NewMovePathTool(filesystem))
	toolRegistry.Register(tools.NewCopyPathTool(filesystem))
	toolRegistry.Register(tools.NewRenamePathTool(filesystem))

	return toolRegistry
}
