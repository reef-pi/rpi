package hal

import (
	"fmt"
	"sort"

	"github.com/reef-pi/embd"

	"github.com/reef-pi/hal"
)

var (
	validGPIOPins = map[int]bool{
		2:  true,
		3:  true,
		4:  true,
		5:  true,
		6:  true,
		7:  true,
		8:  true,
		9:  true,
		10: true,
		11: true,
		12: true,
		13: true,
		14: true,
		15: true,
		16: true,
		17: true,
		18: true,
		19: true,
		20: true,
		21: true,
		22: true,
		23: true,
		24: true,
		25: true,
		26: true,
		27: true,
	}
)

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
	if err := p.digitalPin.SetDirection(embd.In); err != nil {
		return false, fmt.Errorf("can't read input from channel %d: %v", p.number, err)
	}
	v, err := p.digitalPin.Read()
	return v == 1, err
}

func (p *pin) Write(state bool) error {
	if err := p.digitalPin.SetDirection(embd.Out); err != nil {
		return fmt.Errorf("can't set output on channel %d: %v", p.number, err)
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

func (r *driver) DigitalInputPins() []hal.DigitalInputPin {
	var pins []hal.DigitalInputPin
	for _, pin := range r.pins {
		pins = append(pins, pin)
	}
	sort.Slice(pins, func(i, j int) bool { return pins[i].Name() < pins[j].Name() })
	return pins
}

func (r *driver) DigitalInputPin(p int) (hal.DigitalInputPin, error) {
	pin, ok := r.pins[p]
	if !ok {
		return nil, fmt.Errorf("pin %d unknown", p)
	}
	return pin, nil
}

func (r *driver) DigitalOutputPins() []hal.DigitalOutputPin {
	var pins []hal.DigitalOutputPin
	for _, pin := range r.pins {
		pins = append(pins, pin)
	}
	sort.Slice(pins, func(i, j int) bool { return pins[i].Name() < pins[j].Name() })
	return pins
}

func (r *driver) DigitalOutputPin(p int) (hal.DigitalOutputPin, error) {
	pin, ok := r.pins[p]
	if !ok {
		return nil, fmt.Errorf("pin %d unknown", p)
	}
	return pin, nil
}
