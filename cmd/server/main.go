package main

import (
	"context"
	"errors"
	"fmt"
	"os"
)

func main() {
	exitCode, err := runCLI(context.Background(), os.Args[1:], os.Stdout, run)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "flashgate-mcp: %v\n", err)

		if errors.Is(err, errInvalidCLIArguments) {
			os.Exit(2)
		}

		os.Exit(exitCode)
	}

	os.Exit(exitCode)
}
