//go:build windows

package security

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	fileAttributeHidden       = uint32(0x2)
	fileAttributeReparsePoint = uint32(0x400)
)

func isUNCPath(path string) bool {
	if strings.HasPrefix(path, `\\`) || strings.HasPrefix(path, `//`) {
		return true
	}

	volume := filepath.VolumeName(path)
	return strings.HasPrefix(volume, `\\`) || strings.HasPrefix(volume, `//`)
}

func validatePathMetadata(path string, policy Policy) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}

	attributes, err := fileAttributes(info)
	if err != nil {
		return err
	}

	if !policy.AllowHiddenFiles && attributes&fileAttributeHidden != 0 {
		return ErrHiddenPathDenied
	}

	isSymlink := info.Mode()&os.ModeSymlink != 0
	isReparse := attributes&fileAttributeReparsePoint != 0

	if isReparse && !isSymlink {
		return ErrReparsePointDenied
	}

	if isSymlink && !policy.FollowSymlinks {
		return ErrSymlinkDenied
	}

	return nil
}

func fileAttributes(info os.FileInfo) (uint32, error) {
	data, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok || data == nil {
		return 0, nil
	}

	return data.FileAttributes, nil
}

func evalEffectivePath(path string) (string, error) {
	return evalEffectivePathDepth(path, 0)
}

func evalEffectivePathDepth(path string, depth int) (string, error) {
	if depth > 255 {
		return "", errors.New("too many links")
	}

	absolutePath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", err
	}

	volume := filepath.VolumeName(absolutePath)
	remainder := strings.TrimPrefix(absolutePath, volume)
	remainder = strings.TrimLeft(remainder, `\/`)

	current := volume
	if current == "" {
		current = string(filepath.Separator)
	} else if strings.HasPrefix(absolutePath[len(volume):], `\`) || strings.HasPrefix(absolutePath[len(volume):], `/`) {
		current += string(filepath.Separator)
	}

	if remainder == "" {
		return filepath.Clean(current), nil
	}

	parts := strings.Split(remainder, string(filepath.Separator))
	for index := 0; index < len(parts); index++ {
		part := parts[index]
		if part == "" || part == "." {
			continue
		}

		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			return "", err
		}

		if info.Mode()&os.ModeSymlink == 0 {
			continue
		}

		target, err := os.Readlink(current)
		if err != nil {
			return "", err
		}
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(current), target)
		}

		remaining := append([]string{}, parts[index+1:]...)
		nextPath := target
		if len(remaining) > 0 {
			nextPath = filepath.Join(append([]string{target}, remaining...)...)
		}

		resolved, err := evalEffectivePathDepth(nextPath, depth+1)
		if err != nil {
			return "", err
		}

		return resolved, nil
	}

	return filepath.Clean(current), nil
}
