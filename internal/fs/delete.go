package fs

import (
	"os"
	"path/filepath"
)

// Delete deletes a file or directory.
func (f *LocalFileSystem) Delete(path string, recursive bool) error {
	safePath, err := f.guard.ResolveExisting(path)
	if err != nil {
		return err
	}

	info, err := os.Stat(safePath.String())
	if err != nil {
		return err
	}

	if info.IsDir() {
		if recursive {
			if err := f.ensureRecursiveDeleteWithinLimit(safePath.String()); err != nil {
				return err
			}

			return os.RemoveAll(safePath.String())
		}

		entries, err := os.ReadDir(safePath.String())
		if err != nil {
			return err
		}

		if len(entries) > 0 {
			return ErrDirectoryNotEmpty
		}
	}

	return os.Remove(safePath.String())
}

func (f *LocalFileSystem) ensureRecursiveDeleteWithinLimit(root string) error {
	entries := 0

	return filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == root {
			return nil
		}

		entries++
		if entries > f.limits.MaxDeleteEntries {
			return ErrLimitExceeded
		}

		return nil
	})
}
