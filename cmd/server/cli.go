package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/blacksheepkhan/fileserver-mcp/internal/version"
)

var errInvalidCLIArguments = errors.New("invalid CLI arguments")

type serverRunner func(context.Context) error

func runCLI(ctx context.Context, args []string, stdout io.Writer, runServer serverRunner) (int, error) {
	switch len(args) {
	case 0:
		if err := runServer(ctx); err != nil {
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
			return 2, fmt.Errorf("%w: unknown argument: %s\nUse --help for usage.", errInvalidCLIArguments, args[0])
		}
	default:
		return 2, fmt.Errorf("%w: too many arguments\nUse --help for usage.", errInvalidCLIArguments)
	}
}

func helpText() string {
	return `fileserver-mcp

Usage:
  fileserver-mcp
  fileserver-mcp --version
  fileserver-mcp --help

Environment:
  MCP_ROOT    Root directory exposed to MCP clients
`
}
