//go:build !windows
// +build !windows

package hal

import "github.com/warthog618/go-gpiocdev"

const (
	rpiGpioChip = "gpiochip0"
)

type DigitalPin interface {
	SetDirection(string) error
	Read() (int, error)
	Write(int) error
	Close() error
}

func newDigitalPin(i int) (DigitalPin, error) {
	return &digitalPin{}, nil
}

type digitalPin struct {
	pin int
}

func (p *digitalPin) SetDirection(dir string) error {
	return nil
}

func (p *digitalPin) Read() (int, error) {
	in, err := gpiocdev.RequestLine(rpiGpioChip, p.pin, gpiocdev.AsInput)
	if err != nil {
		return 0, err
	}
	return in.Value()
}

func (p *digitalPin) Write(value int) error {
	return nil
}
func (p *digitalPin) Close() error {
	return nil
}
