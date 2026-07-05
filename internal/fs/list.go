package fs

import "os"

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
