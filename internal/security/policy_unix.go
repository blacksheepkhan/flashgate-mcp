//go:build !windows

package security

import (
	"os"
	"path/filepath"
	"strings"
)

func isUNCPath(path string) bool {
	return strings.HasPrefix(path, `\\`) || strings.HasPrefix(path, `//`)
}

func validatePathMetadata(path string, policy Policy) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}

	if !policy.FollowSymlinks && info.Mode()&os.ModeSymlink != 0 {
		return ErrSymlinkDenied
	}

	return nil
}

func evalEffectivePath(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}
