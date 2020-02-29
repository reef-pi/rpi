package hal

import (
	"fmt"

	"github.com/kidoman/embd"
	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/pwm"
)

type DigitalPin interface {
	SetDirection(embd.Direction) error
	Read() (int, error)
	Write(int) error
	Close() error
}

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

type PinFactory func(key interface{}) (DigitalPin, error)

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
