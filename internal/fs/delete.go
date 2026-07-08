package fs

import "os"

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
