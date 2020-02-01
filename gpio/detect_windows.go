// +build windows

package gpio

const (
	bcm2835Base = 0x3F000000
	pi1GPIOBase = 0x3F200000
)

func DetectBase() (int64, error) {
	return int64(pi1GPIOBase + 0x200000), nil
}

func Mmap() ([]uint8, error) {
	return nil, nil
}
