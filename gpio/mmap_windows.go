//go:build windows
// +build windows

package gpio

func munmap(mem8 []uint8) error {
	return nil
}
