package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/thomasweidner/flashgate-mcp/internal/config"
)

func main() {
	os.Exit(runMain(context.Background(), os.Args[1:], os.Stdout, os.Stderr, run))
}

func runMain(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer, runServer serverRunner) int {
	exitCode, err := runCLI(ctx, args, stdout, runServer)
	if err == nil {
		return exitCode
	}

	if errors.Is(err, errInvalidCLIArguments) {
		_, _ = fmt.Fprintf(stderr, "flashgate-mcp: %v\n", err)
		return 2
	}

	category, ok := config.CategoryOf(err)
	if !ok {
		category = config.CategoryStartupFailed
	}
	_, _ = fmt.Fprintf(stderr, "flashgate-mcp: startup failed (%s)\n", category)

	return exitCode
}
