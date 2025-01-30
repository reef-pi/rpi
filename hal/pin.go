package hal

import "fmt"

type DigitalPin interface {
	SetDirection(bool) error
	Read() (int, error)
	Write(int) error
	Close() error
}

type pin struct {
	number    int
	name      string
	lastState bool

	digitalPin DigitalPin
}

func (p *pin) Name() string {
	return p.name
}
func (p *pin) Number() int {
	return p.number
}

func (p *pin) Close() error {
	return p.digitalPin.Close()
}

func (p *pin) Read() (bool, error) {
	if err := p.digitalPin.SetDirection(false); err != nil {
		return false, fmt.Errorf("can't read input from pin %d: %v", p.number, err)
	}
	v, err := p.digitalPin.Read()
	return v == 1, err
}

func (p *pin) Write(state bool) error {
	if err := p.digitalPin.SetDirection(true); err != nil {
		return fmt.Errorf("can't set output on pin %d: %v", p.number, err)
	}
	value := 0
	if state {
		value = 1
	}
	p.lastState = state
	return p.digitalPin.Write(value)
}

func (p *pin) LastState() bool {
	return p.lastState
}
