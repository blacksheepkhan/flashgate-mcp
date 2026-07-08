package fs

import (
	"errors"
	"io"
	"os"
)

// Write writes a file. Existing files are only overwritten when overwrite is true.
func (f *LocalFileSystem) Write(path string, content []byte, overwrite bool) error {
	safePath, err := f.guard.ResolveForCreate(path)
	if err != nil {
		return err
	}

	if int64(len(content)) > f.limits.MaxWriteBytes {
		return ErrLimitExceeded
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

	written, err := file.Write(content)
	if err != nil {
		return err
	}

	if written != len(content) {
		return io.ErrShortWrite
	}

	return nil
}

// Mkdir creates a directory and any missing parent directories.
func (f *LocalFileSystem) Mkdir(path string) error {
	safePath, err := f.guard.ResolveForCreate(path)
	if err != nil {
		return err
	}

	return os.MkdirAll(safePath.String(), 0o700)
}
