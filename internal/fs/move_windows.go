//go:build windows

package fs

import (
	"errors"
	"path/filepath"
	"strings"
	"syscall"
)

func pathsEquivalent(left string, right string) bool {
	return strings.EqualFold(filepath.Clean(left), filepath.Clean(right))
}

func pathsOnSameVolume(source string, target string) (bool, error) {
	return strings.EqualFold(filepath.VolumeName(source), filepath.VolumeName(target)), nil
}

func isCrossVolumeError(err error) bool {
	const errorNotSameDevice syscall.Errno = 17
	return errors.Is(err, errorNotSameDevice)
}
