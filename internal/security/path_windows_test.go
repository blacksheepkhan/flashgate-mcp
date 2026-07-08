//go:build windows

package security

import (
	"errors"
	"path/filepath"
	"syscall"
	"testing"
)

func TestResolveExistingRejectsWindowsHiddenAttributeByDefault(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	path := filepath.Join(root, "hidden.txt")
	writeSecurityTestFile(t, path, "hidden")
	setHiddenAttributeOrSkip(t, path)

	guard := mustNewPathGuard(t, root)

	_, err := guard.ResolveExisting("hidden.txt")

	if !errors.Is(err, ErrHiddenPathDenied) {
		t.Fatalf("expected ErrHiddenPathDenied, got %v", err)
	}
}

func TestResolveExistingAllowsWindowsHiddenAttributeWhenConfigured(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	path := filepath.Join(root, "hidden.txt")
	writeSecurityTestFile(t, path, "hidden")
	setHiddenAttributeOrSkip(t, path)

	guard := mustNewPathGuardWithPolicy(t, root, Policy{AllowHiddenFiles: true})

	_, err := guard.ResolveExisting("hidden.txt")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func setHiddenAttributeOrSkip(t *testing.T, path string) {
	t.Helper()

	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		t.Fatalf("failed to encode path: %v", err)
	}

	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		t.Skipf("failed to read file attributes: %v", err)
	}

	if err := syscall.SetFileAttributes(pointer, attributes|fileAttributeHidden); err != nil {
		t.Skipf("failed to set hidden attribute: %v", err)
	}
}
