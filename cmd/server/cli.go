package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/thomasweidner/flashgate-mcp/internal/config"
	"github.com/thomasweidner/flashgate-mcp/internal/version"
)

var errInvalidCLIArguments = errors.New("invalid CLI arguments")

type serverRunner func(context.Context) error

func runCLI(ctx context.Context, args []string, stdout io.Writer, runServer serverRunner) (int, error) {
	switch len(args) {
	case 0:
		if err := runServer(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				return 0, nil
			}
			if category, ok := config.CategoryOf(err); ok && category != config.CategoryStartupFailed {
				return 3, err
			}
			return 1, err
		}

		return 0, nil
	case 1:
		switch args[0] {
		case "--version":
			_, _ = fmt.Fprintln(stdout, version.Get().String())
			return 0, nil
		case "--help", "-h":
			_, _ = fmt.Fprint(stdout, helpText())
			return 0, nil
		default:
			return 2, fmt.Errorf("%w: unknown argument\nUse --help for usage.", errInvalidCLIArguments)
		}
	default:
		return 2, fmt.Errorf("%w: too many arguments\nUse --help for usage.", errInvalidCLIArguments)
	}
}

func helpText() string {
	return `flashgate-mcp

Usage:
  flashgate-mcp
  flashgate-mcp --version
  flashgate-mcp --help

Environment:
  MCP_ROOT             Required absolute root directory exposed to MCP clients
  MCP_READ_ONLY        Set to true to expose only read-only filesystem tools
  MCP_ALLOW_CWD_ROOT   Development only: set to true with MCP_ROOT=.
`
}
