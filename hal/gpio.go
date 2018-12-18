package hal

import (
	"fmt"
	"sort"

	"github.com/kidoman/embd"

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

	digitalPin embd.DigitalPin
}

func (p *pin) Name() string {
	return p.name
}

func (p *pin) Close() error {
	return p.digitalPin.Close()
}

func (p *pin) Read() (bool, error) {
	err := p.digitalPin.SetDirection(embd.In)
	if err != nil {
		return false, fmt.Errorf("can't read input from channel %d: %v", p.number, err)
	}

	v, err := p.digitalPin.Read()
	if err != nil {
		return false, err
	}
	return v == 1, nil
}

func (p *pin) Write(state bool) error {
	err := p.digitalPin.SetDirection(embd.Out)
	if err != nil {
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

func (r *driver) InputPins() []hal.InputPin {
	var pins []hal.InputPin
	for _, pin := range r.pins {
		pins = append(pins, pin)
	}
	sort.Slice(pins, func(i, j int) bool { return pins[i].Name() < pins[j].Name() })
	return pins
}

func (r *driver) InputPin(name string) (hal.InputPin, error) {
	pin, ok := r.pins[name]
	if !ok {
		return nil, fmt.Errorf("pin %s unknown", name)
	}
	return pin, nil
}

func (r *driver) OutputPins() []hal.OutputPin {
	var pins []hal.OutputPin
	for _, pin := range r.pins {
		pins = append(pins, pin)
	}
	sort.Slice(pins, func(i, j int) bool { return pins[i].Name() < pins[j].Name() })
	return pins
}

func (r *driver) OutputPin(name string) (hal.OutputPin, error) {
	pin, ok := r.pins[name]
	if !ok {
		return nil, fmt.Errorf("pin %s unknown", name)
	}
	return pin, nil
}
