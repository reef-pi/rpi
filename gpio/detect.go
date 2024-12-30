//go:build !windows
// +build !windows

package gpio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
)

const (
	bcm2835Base = 0x3F000000
	pi1GPIOBase = 0x3F200000
)

func DetectBase() (int64, error) {
	base := int64(pi1GPIOBase)
	ranges, err := os.Open("/proc/device-tree/soc/ranges")
	defer ranges.Close()
	if err != nil {
		return base, err
	}
	b := make([]byte, 4)
	n, err := ranges.ReadAt(b, 4)
	if err != nil {
		return base, err
	}
	if n != 4 {
		return base, fmt.Errorf("DT system on chip ranges is %d bytes instead of 4 bytes", n)
	}
	buf := bytes.NewReader(b)
	var out uint32
	err = binary.Read(buf, binary.BigEndian, &out)
	if err != nil {
		return base, nil
	}
	return int64(out + 0x200000), nil
}

func Mmap() ([]uint8, error) {
	memLength := 4096
	var file *os.File
	base, err := DetectBase()
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0)
	switch {
	case os.IsNotExist(err):
		f1, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0)
		if err != nil {
			return nil, err
		}
		file = f1
	case err != nil:
		return nil, err
	default:
		file = f
	}
	defer file.Close()
	// Memory map GPIO registers to byte array
	return syscall.Mmap(
		int(file.Fd()),
		base,
		memLength,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
}
