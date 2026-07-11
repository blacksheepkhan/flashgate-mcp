//go:build !windows

package fs

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
)

func pathsEquivalent(left string, right string) bool {
	return filepath.Clean(left) == filepath.Clean(right)
}

func pathsOnSameVolume(source string, target string) (bool, error) {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return false, err
	}

	targetProbe := target
	for {
		targetInfo, statErr := os.Stat(targetProbe)
		if statErr == nil {
			sourceStat, sourceOK := sourceInfo.Sys().(*syscall.Stat_t)
			targetStat, targetOK := targetInfo.Sys().(*syscall.Stat_t)
			if !sourceOK || !targetOK {
				return true, nil
			}
			return sourceStat.Dev == targetStat.Dev, nil
		}
		if !errors.Is(statErr, os.ErrNotExist) {
			return false, statErr
		}

		parent := filepath.Dir(targetProbe)
		if parent == targetProbe {
			return false, statErr
		}
		targetProbe = parent
	}
}

func isCrossVolumeError(err error) bool {
	return errors.Is(err, syscall.EXDEV)
}
