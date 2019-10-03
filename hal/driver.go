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

type Settings struct {
	PWMFreq int `json:"pwm_freq"`
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

func NewAdapter(s Settings, pd pwm.Driver, factory PinFactory) (*driver, error) {
	d := &driver{
		pins:     make(map[int]*pin),
		channels: make(map[int]*channel),
		meta: hal.Metadata{
			Name:         "rpi",
			Description:  "hardware peripherals and GPIO channels on the base raspberry pi hardware",
			Capabilities: []hal.Capability{hal.DigitalInput, hal.DigitalOutput, hal.PWM},
		},
	}
	for i := range validGPIOPins {
		p, err := factory(i)

		if err != nil {
			return nil, fmt.Errorf("can't build hal pin %d: %v", i, err)
		}
		name := fmt.Sprintf("GP%d", i)
		d.pins[i] = &pin{
			name:       name,
			number:     i,
			digitalPin: p,
		}
	}

	for _, p := range []int{0, 1} {
		ch := &channel{
			pin:       p,
			driver:    pd,
			frequency: s.PWMFreq,
			name:      fmt.Sprintf("%d", p),
		}
		d.channels[p] = ch
	}
	return d, nil
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
