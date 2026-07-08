package fs

import (
	"errors"
	"os"
)

// Move moves a file or directory. Existing targets are only overwritten when overwrite is true.
func (f *LocalFileSystem) Move(source string, target string, overwrite bool) error {
	sourcePath, err := f.guard.ResolveExisting(source)
	if err != nil {
		return err
	}

	targetPath, err := f.guard.ResolveForCreate(target)
	if err != nil {
		return err
	}

	if err := ensureTargetPolicy(targetPath.String(), overwrite); err != nil {
		return err
	}

	if overwrite {
		if err := removeExistingTarget(targetPath.String()); err != nil {
			return err
		}
	}

	return os.Rename(sourcePath.String(), targetPath.String())
}

// Rename renames a file or directory. It is a semantic alias for Move.
func (f *LocalFileSystem) Rename(source string, target string, overwrite bool) error {
	return f.Move(source, target, overwrite)
}

func ensureTargetPolicy(targetPath string, overwrite bool) error {
	info, err := os.Stat(targetPath)
	if err == nil {
		if !overwrite {
			return ErrFileExists
		}

		if info.IsDir() {
			return ErrPathIsDirectory
		}

		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

func removeExistingTarget(targetPath string) error {
	info, err := os.Stat(targetPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	if info.IsDir() {
		return ErrPathIsDirectory
	}

	return os.Remove(targetPath)
}
