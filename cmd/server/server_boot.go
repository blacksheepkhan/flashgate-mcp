package main

import (
	"context"
	"os"

	"github.com/blacksheepkhan/fileserver-mcp/internal/config"
	"github.com/blacksheepkhan/fileserver-mcp/internal/fs"
	"github.com/blacksheepkhan/fileserver-mcp/internal/mcp/server"
	"github.com/blacksheepkhan/fileserver-mcp/internal/security"
)

func run(ctx context.Context) error {
	cfg, err := config.LoadFromEnvironment()
	if err != nil {
		return err
	}

	filesystem, err := fs.NewLocalFileSystemWithPolicy(
		cfg.Filesystem().RootPath(),
		securityPolicyFromConfig(cfg),
	)
	if err != nil {
		return err
	}

	toolRegistry := createToolRegistry(
		filesystem,
		cfg.Filesystem().MaxFileSize(),
		capabilitiesFromReadOnly(cfg.Filesystem().ReadOnly()),
	)
	mcpRouter := createRouter(cfg.Server().Name(), cfg.Server().Version(), toolRegistry)
	mcpServer := server.New(os.Stdin, os.Stdout, mcpRouter)

	return mcpServer.Run(ctx)
}

func securityPolicyFromConfig(cfg config.Config) security.Policy {
	return security.Policy{
		AllowHiddenFiles: cfg.Security().AllowHiddenFiles(),
		AllowUNCPaths:    cfg.Security().AllowUNCPaths(),
		FollowSymlinks:   cfg.Security().FollowSymlinks(),
	}
}
