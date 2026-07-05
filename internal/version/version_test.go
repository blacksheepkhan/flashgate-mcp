package version

import (
	"strings"
	"testing"
)

func TestGetReturnsDefaultBuildMetadata(t *testing.T) {
	info := Get()

	if info.Version != "dev" {
		t.Fatalf("expected default version dev, got %q", info.Version)
	}

	if info.Commit != "unknown" {
		t.Fatalf("expected default commit unknown, got %q", info.Commit)
	}

	if info.Date != "unknown" {
		t.Fatalf("expected default date unknown, got %q", info.Date)
	}
}

func TestInfoString(t *testing.T) {
	info := Info{
		Version: "v1.2.3",
		Commit:  "abc123",
		Date:    "2026-07-05T12:00:00Z",
	}

	output := info.String()

	expectedParts := []string{
		"fileserver-mcp",
		"version: v1.2.3",
		"commit: abc123",
		"date: 2026-07-05T12:00:00Z",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Fatalf("expected output to contain %q, got %q", part, output)
		}
	}
}
