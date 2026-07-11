package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestRunCLIHelp(t *testing.T) {
	var stdout bytes.Buffer

	exitCode, err := runCLI(context.Background(), []string{"--help"}, &stdout, failIfServerRuns(t))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	output := stdout.String()
	expectedParts := []string{
		"flashgate-mcp",
		"Usage:",
		"flashgate-mcp --version",
		"flashgate-mcp --help",
		"MCP_ROOT",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Fatalf("expected help output to contain %q, got %q", part, output)
		}
	}
}

func TestRunCLIShortHelp(t *testing.T) {
	var stdout bytes.Buffer

	exitCode, err := runCLI(context.Background(), []string{"-h"}, &stdout, failIfServerRuns(t))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("expected help output, got %q", stdout.String())
	}
}

func TestRunCLIVersion(t *testing.T) {
	var stdout bytes.Buffer

	exitCode, err := runCLI(context.Background(), []string{"--version"}, &stdout, failIfServerRuns(t))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	output := stdout.String()
	expectedParts := []string{
		"flashgate-mcp",
		"version:",
		"commit:",
		"date:",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Fatalf("expected version output to contain %q, got %q", part, output)
		}
	}
}

func TestRunCLIUnknownArgument(t *testing.T) {
	var stdout bytes.Buffer

	exitCode, err := runCLI(context.Background(), []string{"--unknown"}, &stdout, failIfServerRuns(t))

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, errInvalidCLIArguments) {
		t.Fatalf("expected errInvalidCLIArguments, got %v", err)
	}

	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}

	if !strings.Contains(err.Error(), "unknown argument") {
		t.Fatalf("expected unknown argument error, got %v", err)
	}
	if strings.Contains(err.Error(), "--unknown") {
		t.Fatalf("expected argument value to be omitted, got %v", err)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout output, got %q", stdout.String())
	}
}

func TestRunCLITooManyArguments(t *testing.T) {
	var stdout bytes.Buffer

	exitCode, err := runCLI(context.Background(), []string{"--help", "--version"}, &stdout, failIfServerRuns(t))

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, errInvalidCLIArguments) {
		t.Fatalf("expected errInvalidCLIArguments, got %v", err)
	}

	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}

	if !strings.Contains(err.Error(), "too many arguments") {
		t.Fatalf("expected too many arguments error, got %v", err)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout output, got %q", stdout.String())
	}
}

func TestRunCLIWithoutArgumentsRunsServer(t *testing.T) {
	var stdout bytes.Buffer
	called := false

	exitCode, err := runCLI(
		context.Background(),
		nil,
		&stdout,
		func(context.Context) error {
			called = true
			return nil
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	if !called {
		t.Fatal("expected server runner to be called")
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout output, got %q", stdout.String())
	}
}

func TestRunCLIWithoutArgumentsReturnsServerError(t *testing.T) {
	var stdout bytes.Buffer
	expectedErr := errors.New("server failed")

	exitCode, err := runCLI(
		context.Background(),
		nil,
		&stdout,
		func(context.Context) error {
			return expectedErr
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected server error, got %v", err)
	}

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout output, got %q", stdout.String())
	}
}

func TestRunCLIContextCancellationIsNormalShutdown(t *testing.T) {
	var stdout bytes.Buffer

	exitCode, err := runCLI(
		context.Background(),
		nil,
		&stdout,
		func(context.Context) error { return context.Canceled },
	)

	if err != nil {
		t.Fatalf("expected context cancellation to be treated as normal shutdown, got %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout output, got %q", stdout.String())
	}
}

func failIfServerRuns(t *testing.T) serverRunner {
	t.Helper()

	return func(context.Context) error {
		t.Fatal("server runner must not be called")
		return nil
	}
}
