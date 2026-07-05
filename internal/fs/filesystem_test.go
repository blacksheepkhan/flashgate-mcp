package fs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/blacksheepkhan/fileserver-mcp/internal/security"
)

func TestNewLocalFileSystemRejectsEmptyRoot(t *testing.T) {
	t.Parallel()

	_, err := NewLocalFileSystem("")

	if !errors.Is(err, security.ErrEmptyRoot) {
		t.Fatalf("expected ErrEmptyRoot, got %v", err)
	}
}

func TestLocalFileSystemListReturnsFilesAndDirectories(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	writeTestFile(t, filepath.Join(root, "file.txt"), "hello")
	mkdir(t, filepath.Join(root, "subdir"))

	filesystem := mustNewLocalFileSystem(t, root)

	entries, err := filesystem.List(".")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d: %#v", len(entries), entries)
	}

	fileEntry := findEntry(t, entries, "file.txt")
	if fileEntry.IsDir {
		t.Fatal("expected file.txt to be a file")
	}

	if fileEntry.Size != int64(len("hello")) {
		t.Fatalf("expected file size %d, got %d", len("hello"), fileEntry.Size)
	}

	dirEntry := findEntry(t, entries, "subdir")
	if !dirEntry.IsDir {
		t.Fatal("expected subdir to be a directory")
	}
}

func TestLocalFileSystemListEmptyDirectory(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	entries, err := filesystem.List(".")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(entries) != 0 {
		t.Fatalf("expected empty directory, got %d entries", len(entries))
	}
}

func TestLocalFileSystemListNestedDirectory(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	nestedDir := filepath.Join(root, "alpha", "beta")

	mkdir(t, nestedDir)
	writeTestFile(t, filepath.Join(nestedDir, "nested.txt"), "nested-content")

	filesystem := mustNewLocalFileSystem(t, root)

	entries, err := filesystem.List(filepath.Join("alpha", "beta"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	entry := findEntry(t, entries, "nested.txt")
	if entry.IsDir {
		t.Fatal("expected nested.txt to be a file")
	}

	if entry.Size != int64(len("nested-content")) {
		t.Fatalf("expected size %d, got %d", len("nested-content"), entry.Size)
	}
}

func TestLocalFileSystemListRejectsPathTraversal(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	_, err := filesystem.List("..")

	if !errors.Is(err, security.ErrPathTraversal) {
		t.Fatalf("expected ErrPathTraversal, got %v", err)
	}
}

func TestLocalFileSystemListRejectsAbsolutePath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	_, err := filesystem.List(filepath.Join(root, "file.txt"))

	if !errors.Is(err, security.ErrAbsolutePath) {
		t.Fatalf("expected ErrAbsolutePath, got %v", err)
	}
}

func TestLocalFileSystemListReturnsErrorForMissingDirectory(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	_, err := filesystem.List("missing")

	if err == nil {
		t.Fatal("expected error for missing directory")
	}
}

func TestLocalFileSystemReadReturnsFileContent(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "file.txt"), "hello world")

	filesystem := mustNewLocalFileSystem(t, root)

	content, err := filesystem.Read("file.txt", 1024)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(content) != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", string(content))
	}
}

func TestLocalFileSystemReadRejectsDirectory(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mkdir(t, filepath.Join(root, "subdir"))

	filesystem := mustNewLocalFileSystem(t, root)

	_, err := filesystem.Read("subdir", 1024)

	if !errors.Is(err, ErrPathIsDirectory) {
		t.Fatalf("expected ErrPathIsDirectory, got %v", err)
	}
}

func TestLocalFileSystemReadRejectsTooLargeFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "file.txt"), "hello world")

	filesystem := mustNewLocalFileSystem(t, root)

	_, err := filesystem.Read("file.txt", 5)

	if !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("expected ErrFileTooLarge, got %v", err)
	}
}

func TestLocalFileSystemReadRejectsZeroLimit(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "file.txt"), "hello")

	filesystem := mustNewLocalFileSystem(t, root)

	_, err := filesystem.Read("file.txt", 0)

	if !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("expected ErrFileTooLarge, got %v", err)
	}
}

func TestLocalFileSystemReadRejectsTraversal(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	_, err := filesystem.Read("..", 1024)

	if !errors.Is(err, security.ErrPathTraversal) {
		t.Fatalf("expected ErrPathTraversal, got %v", err)
	}
}

func TestLocalFileSystemStatReturnsMetadata(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "file.txt"), "hello")

	filesystem := mustNewLocalFileSystem(t, root)

	metadata, err := filesystem.Stat("file.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if metadata.Name != "file.txt" {
		t.Fatalf("expected name %q, got %q", "file.txt", metadata.Name)
	}

	if metadata.IsDir {
		t.Fatal("expected file.txt to be a file")
	}

	if metadata.Size != int64(len("hello")) {
		t.Fatalf("expected size %d, got %d", len("hello"), metadata.Size)
	}
}

func TestLocalFileSystemStatReturnsDirectoryMetadata(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mkdir(t, filepath.Join(root, "subdir"))

	filesystem := mustNewLocalFileSystem(t, root)

	metadata, err := filesystem.Stat("subdir")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if metadata.Name != "subdir" {
		t.Fatalf("expected name %q, got %q", "subdir", metadata.Name)
	}

	if !metadata.IsDir {
		t.Fatal("expected subdir to be a directory")
	}
}

func TestLocalFileSystemStatRejectsTraversal(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	_, err := filesystem.Stat("..")

	if !errors.Is(err, security.ErrPathTraversal) {
		t.Fatalf("expected ErrPathTraversal, got %v", err)
	}
}

func TestLocalFileSystemExistsReturnsTrueForExistingFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "file.txt"), "hello")

	filesystem := mustNewLocalFileSystem(t, root)

	exists, err := filesystem.Exists("file.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !exists {
		t.Fatal("expected file to exist")
	}
}

func TestLocalFileSystemExistsReturnsFalseForMissingFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	exists, err := filesystem.Exists("missing.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if exists {
		t.Fatal("expected file to be missing")
	}
}

func TestLocalFileSystemExistsRejectsTraversal(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	exists, err := filesystem.Exists("..")

	if !errors.Is(err, security.ErrPathTraversal) {
		t.Fatalf("expected ErrPathTraversal, got %v", err)
	}

	if exists {
		t.Fatal("expected traversal path to not exist")
	}
}

func TestLocalFileSystemWriteCreatesFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	err := filesystem.Write("created.txt", []byte("created-content"), false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	content := readTestFile(t, filepath.Join(root, "created.txt"))
	if content != "created-content" {
		t.Fatalf("expected %q, got %q", "created-content", content)
	}
}

func TestLocalFileSystemWriteRejectsExistingFileWithoutOverwrite(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "file.txt"), "old-content")

	filesystem := mustNewLocalFileSystem(t, root)

	err := filesystem.Write("file.txt", []byte("new-content"), false)

	if !errors.Is(err, ErrFileExists) {
		t.Fatalf("expected ErrFileExists, got %v", err)
	}

	content := readTestFile(t, filepath.Join(root, "file.txt"))
	if content != "old-content" {
		t.Fatalf("expected file content to remain unchanged, got %q", content)
	}
}

func TestLocalFileSystemWriteOverwritesExistingFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "file.txt"), "old-content")

	filesystem := mustNewLocalFileSystem(t, root)

	err := filesystem.Write("file.txt", []byte("new-content"), true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	content := readTestFile(t, filepath.Join(root, "file.txt"))
	if content != "new-content" {
		t.Fatalf("expected %q, got %q", "new-content", content)
	}
}

func TestLocalFileSystemWriteRejectsDirectory(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mkdir(t, filepath.Join(root, "subdir"))

	filesystem := mustNewLocalFileSystem(t, root)

	err := filesystem.Write("subdir", []byte("content"), true)

	if !errors.Is(err, ErrPathIsDirectory) {
		t.Fatalf("expected ErrPathIsDirectory, got %v", err)
	}
}

func TestLocalFileSystemWriteRejectsTraversal(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	err := filesystem.Write("..", []byte("content"), false)

	if !errors.Is(err, security.ErrPathTraversal) {
		t.Fatalf("expected ErrPathTraversal, got %v", err)
	}
}

func TestLocalFileSystemWriteRejectsAbsolutePath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	err := filesystem.Write(filepath.Join(root, "file.txt"), []byte("content"), false)

	if !errors.Is(err, security.ErrAbsolutePath) {
		t.Fatalf("expected ErrAbsolutePath, got %v", err)
	}
}

func TestLocalFileSystemMkdirCreatesDirectory(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	err := filesystem.Mkdir(filepath.Join("alpha", "beta"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	info, err := os.Stat(filepath.Join(root, "alpha", "beta"))
	if err != nil {
		t.Fatalf("expected directory to exist, got %v", err)
	}

	if !info.IsDir() {
		t.Fatal("expected created path to be a directory")
	}
}

func TestLocalFileSystemMkdirRejectsTraversal(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	err := filesystem.Mkdir("..")

	if !errors.Is(err, security.ErrPathTraversal) {
		t.Fatalf("expected ErrPathTraversal, got %v", err)
	}
}

func TestLocalFileSystemMkdirRejectsAbsolutePath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filesystem := mustNewLocalFileSystem(t, root)

	err := filesystem.Mkdir(filepath.Join(root, "alpha"))

	if !errors.Is(err, security.ErrAbsolutePath) {
		t.Fatalf("expected ErrAbsolutePath, got %v", err)
	}
}

func mustNewLocalFileSystem(t *testing.T, root string) *LocalFileSystem {
	t.Helper()

	filesystem, err := NewLocalFileSystem(root)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	return filesystem
}

func writeTestFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	return string(content)
}

func mkdir(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}
}

func findEntry(t *testing.T, entries []Entry, name string) Entry {
	t.Helper()

	for _, entry := range entries {
		if entry.Name == name {
			return entry
		}
	}

	t.Fatalf("entry %q not found in %#v", name, entries)

	return Entry{}
}
