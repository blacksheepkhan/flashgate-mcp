package main

import (
	"context"
	"os"

	"github.com/blacksheepkhan/fileserver-mcp/internal/config"
	"github.com/blacksheepkhan/fileserver-mcp/internal/fs"
	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/initialize"
	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/router"
	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/server"
	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/tools"
)

func run(ctx context.Context) error {
	cfg, err := config.LoadFromEnvironment()
	if err != nil {
		return err
	}

	filesystem, err := fs.NewLocalFileSystem(cfg.Filesystem().RootPath())
	if err != nil {
		return err
	}

	toolRegistry := createToolRegistry(filesystem, cfg.Filesystem().MaxFileSize())

	mcpRouter := router.New()
	mcpRouter.Register(initialize.NewHandler(cfg.Server().Name(), cfg.Server().Version()))
	mcpRouter.Register(tools.NewListHandler(toolRegistry))
	mcpRouter.Register(tools.NewCallHandler(toolRegistry))

	mcpServer := server.New(os.Stdin, os.Stdout, mcpRouter)

	return mcpServer.Run(ctx)
}
