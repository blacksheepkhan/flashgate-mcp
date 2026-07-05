package fs

import (
	"errors"

	"github.com/blacksheepkhan/fileserver-mcp/internal/security"
)

var (
	// ErrFileTooLarge is returned when a file exceeds the configured read limit.
	ErrFileTooLarge = errors.New("file exceeds maximum allowed size")

	// ErrPathIsDirectory is returned when a file operation receives a directory path.
	ErrPathIsDirectory = errors.New("path is a directory")

	// ErrPathIsNotDirectory is returned when a directory operation receives a non-directory path.
	ErrPathIsNotDirectory = errors.New("path is not a directory")

	// ErrFileExists is returned when writing a file that already exists without overwrite permission.
	ErrFileExists = errors.New("file already exists")

	// ErrDirectoryNotEmpty is returned when deleting a non-empty directory without recursive deletion.
	ErrDirectoryNotEmpty = errors.New("directory is not empty")

	// ErrCopyDirectoryUnsupported is returned when attempting to copy a directory.
	ErrCopyDirectoryUnsupported = errors.New("copying directories is not supported")
)

// Entry represents a filesystem directory entry.
type Entry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

// Metadata represents filesystem metadata.
type Metadata struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

// FileSystem defines filesystem operations used by MCP tools.
type FileSystem interface {
	List(path string) ([]Entry, error)
	Read(path string, maxBytes int64) ([]byte, error)
	Stat(path string) (Metadata, error)
	Exists(path string) (bool, error)
	Write(path string, content []byte, overwrite bool) error
	Mkdir(path string) error
	Delete(path string, recursive bool) error
	Move(source string, target string, overwrite bool) error
	Copy(source string, target string, overwrite bool) error
	Rename(source string, target string, overwrite bool) error
}

// LocalFileSystem implements FileSystem using the local operating system.
type LocalFileSystem struct {
	guard *security.PathGuard
}

// NewLocalFileSystem creates a new local filesystem.
func NewLocalFileSystem(root string) (*LocalFileSystem, error) {
	guard, err := security.NewPathGuard(root)
	if err != nil {
		return nil, err
	}

	return &LocalFileSystem{
		guard: guard,
	}, nil
}
