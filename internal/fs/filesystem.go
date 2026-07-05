package fs

import (
	"errors"
	"os"

	"github.com/blacksheepkhan/fileserver-mcp/internal/security"
)

var (
	// ErrFileTooLarge is returned when a file exceeds the configured read limit.
	ErrFileTooLarge = errors.New("file exceeds maximum allowed size")

	// ErrPathIsDirectory is returned when a file operation receives a directory path.
	ErrPathIsDirectory = errors.New("path is a directory")

	// ErrFileExists is returned when writing a file that already exists without overwrite permission.
	ErrFileExists = errors.New("file already exists")
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

// Read reads a file up to maxBytes bytes.
func (f *LocalFileSystem) Read(path string, maxBytes int64) ([]byte, error) {
	safePath, err := f.guard.Resolve(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(safePath.String())
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, ErrPathIsDirectory
	}

	if maxBytes <= 0 {
		return nil, ErrFileTooLarge
	}

	if info.Size() > maxBytes {
		return nil, ErrFileTooLarge
	}

	return os.ReadFile(safePath.String())
}

// Stat returns filesystem metadata.
func (f *LocalFileSystem) Stat(path string) (Metadata, error) {
	safePath, err := f.guard.Resolve(path)
	if err != nil {
		return Metadata{}, err
	}

	info, err := os.Stat(safePath.String())
	if err != nil {
		return Metadata{}, err
	}

	return Metadata{
		Name:  info.Name(),
		IsDir: info.IsDir(),
		Size:  info.Size(),
	}, nil
}

// Exists checks whether a path exists.
func (f *LocalFileSystem) Exists(path string) (bool, error) {
	safePath, err := f.guard.Resolve(path)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(safePath.String())
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

// Write writes a file. Existing files are only overwritten when overwrite is true.
func (f *LocalFileSystem) Write(path string, content []byte, overwrite bool) error {
	safePath, err := f.guard.Resolve(path)
	if err != nil {
		return err
	}

	info, err := os.Stat(safePath.String())
	if err == nil {
		if info.IsDir() {
			return ErrPathIsDirectory
		}

		if !overwrite {
			return ErrFileExists
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	flags := os.O_WRONLY | os.O_CREATE
	if overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	file, err := os.OpenFile(safePath.String(), flags, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return ErrFileExists
		}

		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	return err
}

// Mkdir creates a directory and any missing parent directories.
func (f *LocalFileSystem) Mkdir(path string) error {
	safePath, err := f.guard.Resolve(path)
	if err != nil {
		return err
	}

	return os.MkdirAll(safePath.String(), 0o700)
}
