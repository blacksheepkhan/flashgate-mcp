//go:build windows

package fs

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestWindowsMoveRejectsCaseAlias(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "Case.txt")
	writeTestFile(t, path, "content")
	filesystem := mustNewLocalFileSystem(t, root)

	err := filesystem.Move("Case.txt", "case.txt", true)
	if !errors.Is(err, ErrSamePath) {
		t.Fatalf("expected ErrSamePath, got %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected source to remain: %v", err)
	}
}

func TestWindowsCrossVolumeErrorClassification(t *testing.T) {
	if !isCrossVolumeError(syscall.Errno(17)) {
		t.Fatal("expected ERROR_NOT_SAME_DEVICE to be classified")
	}
}
