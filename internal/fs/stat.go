package fs

import (
	"errors"
	"os"
)

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
