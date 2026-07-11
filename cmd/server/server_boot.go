package main

import (
	"context"
	"os"

	"github.com/blacksheepkhan/flashgate-mcp/internal/config"
	"github.com/blacksheepkhan/flashgate-mcp/internal/diagnostics"
	"github.com/blacksheepkhan/flashgate-mcp/internal/fs"
	"github.com/blacksheepkhan/flashgate-mcp/internal/mcp/server"
	"github.com/blacksheepkhan/flashgate-mcp/internal/security"
)

func run(ctx context.Context) error {
	cfg, err := config.LoadFromEnvironment()
	if err != nil {
		return err
	}

	filesystem, err := fs.NewLocalFileSystemWithPolicyAndLimits(
		cfg.Filesystem().RootPath(),
		securityPolicyFromConfig(cfg),
		filesystemLimitsFromConfig(cfg),
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
	mcpServer := server.NewWithOptions(os.Stdin, os.Stdout, mcpRouter, serverOptionsFromConfig(cfg))

	return mcpServer.Run(ctx)
}

func filesystemLimitsFromConfig(cfg config.Config) fs.Limits {
	return fs.Limits{
		MaxWriteBytes:    cfg.Filesystem().MaxWriteBytes(),
		MaxListEntries:   cfg.Filesystem().MaxListEntries(),
		MaxCopyBytes:     cfg.Filesystem().MaxCopyBytes(),
		MaxDeleteEntries: cfg.Filesystem().MaxDeleteEntries(),
	}
}

func securityPolicyFromConfig(cfg config.Config) security.Policy {
	return security.Policy{
		AllowHiddenFiles: cfg.Security().AllowHiddenFiles(),
		AllowUNCPaths:    cfg.Security().AllowUNCPaths(),
		FollowSymlinks:   cfg.Security().FollowSymlinks(),
	}
}

func serverOptionsFromConfig(cfg config.Config) server.Options {
	return server.Options{
		MaxMessageBytes:  cfg.Server().MaxMessageBytes(),
		MaxArgumentBytes: cfg.Server().MaxArgumentBytes(),
		MaxResponseBytes: cfg.Server().MaxResponseBytes(),
		Diagnostics:      diagnostics.NewLogger(cfg.Server().Debug(), os.Stderr),
	}
}
