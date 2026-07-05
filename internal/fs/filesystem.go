package fs

import (
	"os"

	"github.com/blacksheepkhan/fileserver-mcp/internal/security"
)

// Entry represents a filesystem directory entry.
type Entry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

// FileSystem defines filesystem operations used by MCP tools.
type FileSystem interface {
	List(path string) ([]Entry, error)
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

// List lists directory entries.
func (f *LocalFileSystem) List(path string) ([]Entry, error) {
	safePath, err := f.guard.Resolve(path)
	if err != nil {
		return nil, err
	}

	dirEntries, err := os.ReadDir(safePath.String())
	if err != nil {
		return nil, err
	}

	result := make([]Entry, 0, len(dirEntries))
	for _, dirEntry := range dirEntries {
		info, err := dirEntry.Info()
		if err != nil {
			return nil, err
		}

		result = append(result, Entry{
			Name:  dirEntry.Name(),
			IsDir: dirEntry.IsDir(),
			Size:  info.Size(),
		})
	}

	return result, nil
}
