package main

import (
	"testing"

	"github.com/blacksheepkhan/fileserver-mcp/internal/config"
)

func TestSecurityPolicyFromConfigUsesSecureDefaults(t *testing.T) {
	t.Parallel()

	policy := securityPolicyFromConfig(config.DefaultConfig())

	if policy.AllowHiddenFiles {
		t.Fatal("expected hidden files to be denied by default")
	}

	if policy.AllowUNCPaths {
		t.Fatal("expected UNC paths to be denied by default")
	}

	if policy.FollowSymlinks {
		t.Fatal("expected symlink following to be denied by default")
	}
}

func TestSecurityPolicyFromConfigUsesEnvironment(t *testing.T) {
	t.Setenv("MCP_ALLOW_HIDDEN_FILES", "true")
	t.Setenv("MCP_ALLOW_UNC_PATHS", "true")
	t.Setenv("MCP_FOLLOW_SYMLINKS", "true")

	cfg, err := config.LoadFromEnvironment()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	policy := securityPolicyFromConfig(cfg)

	if !policy.AllowHiddenFiles {
		t.Fatal("expected hidden files to be allowed")
	}

	if !policy.AllowUNCPaths {
		t.Fatal("expected UNC paths to be allowed")
	}

	if !policy.FollowSymlinks {
		t.Fatal("expected symlink following to be allowed")
	}
}
