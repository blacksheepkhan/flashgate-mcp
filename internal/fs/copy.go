package fs

import (
	"errors"
	"io"
	"os"
)

// Copy copies a file. Directory copy is intentionally unsupported.
func (f *LocalFileSystem) Copy(source string, target string, overwrite bool) error {
	sourcePath, err := f.guard.ResolveExisting(source)
	if err != nil {
		return err
	}

	targetPath, err := f.guard.ResolveForCreate(target)
	if err != nil {
		return err
	}

	sourceInfo, err := os.Stat(sourcePath.String())
	if err != nil {
		return err
	}

	if sourceInfo.IsDir() {
		return ErrCopyDirectoryUnsupported
	}

	if sourceInfo.Size() > f.limits.MaxCopyBytes {
		return ErrLimitExceeded
	}

	if err := ensureTargetPolicy(targetPath.String(), overwrite); err != nil {
		return err
	}

	sourceFile, err := os.Open(sourcePath.String())
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	flags := os.O_WRONLY | os.O_CREATE
	if overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	targetFile, err := os.OpenFile(targetPath.String(), flags, sourceInfo.Mode().Perm())
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return ErrFileExists
		}

		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)
	return err
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
