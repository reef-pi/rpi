//go:build !windows
// +build !windows

package gpio

import "syscall"

func munmap(mem8 []uint8) error {
	return syscall.Munmap(mem8)
}
