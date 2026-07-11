//go:build !windows

package fs

import (
	"syscall"
	"testing"
)

func TestUnixCrossVolumeErrorClassification(t *testing.T) {
	if !isCrossVolumeError(syscall.EXDEV) {
		t.Fatal("expected EXDEV to be classified")
	}
}
