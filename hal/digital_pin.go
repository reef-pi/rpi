//go:build !darwin

package hal

import (
	"github.com/warthog618/go-gpiocdev"
)

const (
	rpiGpioChip = "gpiochip0"
)

func newDigitalPin(i int) (DigitalPin, error) {
	return &digitalPin{pin: i}, nil
}

type digitalPin struct {
	pin int
}

func (p *digitalPin) SetDirection(dir bool) error {
	var pinDirection gpiocdev.LineConfigOption = gpiocdev.AsInput
	if dir {
		pinDirection = gpiocdev.AsOutput(1, 0)
	}
	l, err := gpiocdev.RequestLine(rpiGpioChip, p.pin)
	if err != nil {
		return err
	}
	defer l.Close()
	return l.Reconfigure(pinDirection)
}

func (p *digitalPin) Read() (int, error) {
	in, err := gpiocdev.RequestLine(rpiGpioChip, p.pin, gpiocdev.AsInput)
	if err != nil {
		return 0, err
	}
	defer in.Close()
	return in.Value()
}

func (p *digitalPin) Write(value int) error {
	if err := p.SetDirection(true); err != nil {
		return err
	}
	out, err := gpiocdev.RequestLine(rpiGpioChip, p.pin, gpiocdev.AsOutput(value, 0))
	if err != nil {
		defer out.Close()
	}
	return err
}

func (p *digitalPin) Close() error {
	return nil
}
