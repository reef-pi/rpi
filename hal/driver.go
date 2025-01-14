package hal

import (
	"fmt"
	"sort"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/pwm"
)

type driver struct {
	meta      hal.Metadata
	pins      map[int]*pin
	channels  map[int]*channel
	pwmDriver pwm.Driver
}

func (r *driver) Metadata() hal.Metadata {
	return r.meta
}

func (r *driver) Close() error {
	for _, p := range r.pins {
		err := p.Close()
		if err != nil {
			return fmt.Errorf("can't close hal driver due to channel %s", p.Name())
		}
	}
	return nil
}

func (d *driver) Pins(cap hal.Capability) ([]hal.Pin, error) {
	var pins []hal.Pin
	switch cap {
	case hal.DigitalInput, hal.DigitalOutput:
		for _, pin := range d.pins {
			pins = append(pins, pin)
		}
		return pins, nil
	case hal.PWM:
		for _, pin := range d.channels {
			pins = append(pins, pin)
		}
		return pins, nil
	default:
		return nil, fmt.Errorf("Unsupported capability:%s", cap.String())
	}
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
