//go:build !linux

package hal

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
	return nil
}

func (p *digitalPin) Read() (int, error) {
	return 0, nil
}

func (p *digitalPin) Write(value int) error {
	return nil
}

func (p *digitalPin) Close() error {
	return nil
}
