package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()

	if cfg.Filesystem().RootPath() != defaultRootPath {
		t.Fatalf("expected root path %q, got %q", defaultRootPath, cfg.Filesystem().RootPath())
	}

	if cfg.Filesystem().ReadOnly() {
		t.Fatal("expected read-only mode to be disabled by default")
	}

	if cfg.Filesystem().MaxFileSize() != defaultMaxFileSize {
		t.Fatalf("expected max file size %d, got %d", defaultMaxFileSize, cfg.Filesystem().MaxFileSize())
	}

	if cfg.Security().AllowHiddenFiles() {
		t.Fatal("expected hidden files to be denied by default")
	}

	if cfg.Security().AllowUNCPaths() {
		t.Fatal("expected UNC paths to be denied by default")
	}

	if cfg.Security().FollowSymlinks() {
		t.Fatal("expected symlink following to be disabled by default")
	}

	if cfg.Server().Name() != defaultServerName {
		t.Fatalf("expected server name %q, got %q", defaultServerName, cfg.Server().Name())
	}

	if cfg.Server().Version() != defaultVersion {
		t.Fatalf("expected version %q, got %q", defaultVersion, cfg.Server().Version())
	}

	if cfg.Server().Debug() {
		t.Fatal("expected debug mode to be disabled by default")
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	t.Setenv(envRootPath, `C:\temp\fileserver-mcp`)
	t.Setenv(envReadOnly, "true")
	t.Setenv(envMaxFileSize, "1048576")
	t.Setenv(envAllowHiddenFiles, "true")
	t.Setenv(envAllowUNCPaths, "true")
	t.Setenv(envFollowSymlinks, "true")
	t.Setenv(envServerDebug, "true")

	cfg, err := LoadFromEnvironment()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Filesystem().RootPath() != `C:\temp\fileserver-mcp` {
		t.Fatalf("unexpected root path: %q", cfg.Filesystem().RootPath())
	}

	if !cfg.Filesystem().ReadOnly() {
		t.Fatal("expected read-only mode to be enabled")
	}

	if cfg.Filesystem().MaxFileSize() != 1048576 {
		t.Fatalf("expected max file size 1048576, got %d", cfg.Filesystem().MaxFileSize())
	}

	if !cfg.Security().AllowHiddenFiles() {
		t.Fatal("expected hidden files to be allowed")
	}

	if !cfg.Security().AllowUNCPaths() {
		t.Fatal("expected UNC paths to be allowed")
	}

	if !cfg.Security().FollowSymlinks() {
		t.Fatal("expected symlink following to be enabled")
	}

	if !cfg.Server().Debug() {
		t.Fatal("expected debug mode to be enabled")
	}
}

func TestLoadFromEnvironmentParsesFalseSecurityFlags(t *testing.T) {
	t.Setenv(envAllowHiddenFiles, "false")
	t.Setenv(envAllowUNCPaths, "false")
	t.Setenv(envFollowSymlinks, "false")

	cfg, err := LoadFromEnvironment()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Security().AllowHiddenFiles() {
		t.Fatal("expected hidden files to be denied")
	}

	if cfg.Security().AllowUNCPaths() {
		t.Fatal("expected UNC paths to be denied")
	}

	if cfg.Security().FollowSymlinks() {
		t.Fatal("expected symlink following to be denied")
	}
}

func TestLoadFromEnvironmentRejectsInvalidReadOnly(t *testing.T) {
	t.Setenv(envReadOnly, "not-a-bool")

	_, err := LoadFromEnvironment()
	if err == nil {
		t.Fatal("expected error for invalid read-only value")
	}
}

func TestLoadFromEnvironmentRejectsInvalidMaxFileSize(t *testing.T) {
	t.Setenv(envMaxFileSize, "not-a-number")

	_, err := LoadFromEnvironment()
	if err == nil {
		t.Fatal("expected error for invalid max file size value")
	}
}

func TestLoadFromEnvironmentRejectsInvalidAllowHiddenFiles(t *testing.T) {
	t.Setenv(envAllowHiddenFiles, "not-a-bool")

	_, err := LoadFromEnvironment()
	if err == nil {
		t.Fatal("expected error for invalid hidden files value")
	}
}

func TestLoadFromEnvironmentRejectsInvalidAllowUNCPaths(t *testing.T) {
	t.Setenv(envAllowUNCPaths, "not-a-bool")

	_, err := LoadFromEnvironment()
	if err == nil {
		t.Fatal("expected error for invalid UNC paths value")
	}
}

func TestLoadFromEnvironmentRejectsInvalidFollowSymlinks(t *testing.T) {
	t.Setenv(envFollowSymlinks, "not-a-bool")

	_, err := LoadFromEnvironment()
	if err == nil {
		t.Fatal("expected error for invalid symlink value")
	}
}

func TestLoadFromEnvironmentRejectsInvalidDebug(t *testing.T) {
	t.Setenv(envServerDebug, "not-a-bool")

	_, err := LoadFromEnvironment()
	if err == nil {
		t.Fatal("expected error for invalid debug value")
	}
}

func TestValidateAcceptsDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected default config to be valid, got %v", err)
	}
}

func TestValidateRejectsEmptyRootPath(t *testing.T) {
	t.Parallel()

	cfg := Config{
		filesystem: FilesystemConfig{
			rootPath:    "",
			readOnly:    false,
			maxFileSize: defaultMaxFileSize,
		},
		security: SecurityConfig{},
		server:   ServerConfig{},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for empty root path")
	}
}

func TestValidateRejectsZeroMaxFileSize(t *testing.T) {
	t.Parallel()

	cfg := Config{
		filesystem: FilesystemConfig{
			rootPath:    defaultRootPath,
			readOnly:    false,
			maxFileSize: 0,
		},
		security: SecurityConfig{},
		server:   ServerConfig{},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for zero max file size")
	}
}

func TestValidateRejectsNegativeMaxFileSize(t *testing.T) {
	t.Parallel()

	cfg := Config{
		filesystem: FilesystemConfig{
			rootPath:    defaultRootPath,
			readOnly:    false,
			maxFileSize: -1,
		},
		security: SecurityConfig{},
		server:   ServerConfig{},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for negative max file size")
	}
}
