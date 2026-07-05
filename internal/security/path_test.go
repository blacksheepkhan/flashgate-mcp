package security

import (
	"errors"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewPathGuardRejectsEmptyRoot(t *testing.T) {
	t.Parallel()

	_, err := NewPathGuard("")

	if !errors.Is(err, ErrEmptyRoot) {
		t.Fatalf("expected ErrEmptyRoot, got %v", err)
	}
}

func TestNewPathGuardNormalizesRootToAbsolutePath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	guard, err := NewPathGuard(root)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !filepath.IsAbs(guard.Root()) {
		t.Fatalf("expected absolute root path, got %q", guard.Root())
	}

	expected, err := filepath.Abs(filepath.Clean(root))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if guard.Root() != expected {
		t.Fatalf("expected root %q, got %q", expected, guard.Root())
	}
}

func TestResolveAcceptsEmptyUserPathAsRoot(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	guard := mustNewPathGuard(t, root)

	safePath, err := guard.Resolve("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected, err := filepath.Abs(root)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if safePath.String() != expected {
		t.Fatalf("expected %q, got %q", expected, safePath.String())
	}
}

func TestResolveAcceptsRelativePath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	guard := mustNewPathGuard(t, root)

	safePath, err := guard.Resolve(filepath.Join("alpha", "beta.txt"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := filepath.Join(guard.Root(), "alpha", "beta.txt")

	if safePath.String() != expected {
		t.Fatalf("expected %q, got %q", expected, safePath.String())
	}
}

func TestResolveRejectsParentTraversal(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	guard := mustNewPathGuard(t, root)

	_, err := guard.Resolve("..")

	if !errors.Is(err, ErrPathTraversal) {
		t.Fatalf("expected ErrPathTraversal, got %v", err)
	}
}

func TestResolveRejectsNestedParentTraversal(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	guard := mustNewPathGuard(t, root)

	_, err := guard.Resolve(filepath.Join("..", "outside.txt"))

	if !errors.Is(err, ErrPathTraversal) {
		t.Fatalf("expected ErrPathTraversal, got %v", err)
	}
}

func TestResolveRejectsAbsolutePath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	guard := mustNewPathGuard(t, root)

	absolutePath := filepath.Join(root, "file.txt")
	_, err := guard.Resolve(absolutePath)

	if !errors.Is(err, ErrAbsolutePath) {
		t.Fatalf("expected ErrAbsolutePath, got %v", err)
	}
}

func TestSafePathBaseAndDir(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	guard := mustNewPathGuard(t, root)

	safePath, err := guard.Resolve(filepath.Join("dir", "file.txt"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if safePath.Base() != "file.txt" {
		t.Fatalf("expected base %q, got %q", "file.txt", safePath.Base())
	}

	expectedDir := filepath.Join(guard.Root(), "dir")
	if safePath.Dir() != expectedDir {
		t.Fatalf("expected dir %q, got %q", expectedDir, safePath.Dir())
	}
}

func TestResolveRejectsWindowsAbsolutePathOnWindows(t *testing.T) {
	t.Parallel()

	if runtime.GOOS != "windows" {
		t.Skip("windows-specific test")
	}

	root := t.TempDir()
	guard := mustNewPathGuard(t, root)

	_, err := guard.Resolve(`C:\Windows\System32`)

	if !errors.Is(err, ErrAbsolutePath) {
		t.Fatalf("expected ErrAbsolutePath, got %v", err)
	}
}

func mustNewPathGuard(t *testing.T, root string) *PathGuard {
	t.Helper()

	guard, err := NewPathGuard(root)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	return guard
}
