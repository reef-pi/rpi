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
	PWMFreq int
}

type driver struct {
	meta      hal.Metadata
	pins      map[string]*pin
	channels  map[string]*channel
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

//embd.NewDigitalPin
type PinFactory func(key interface{}) (DigitalPin, error)

func NewAdapter(s Settings, pd pwm.Driver, factory PinFactory) (*driver, error) {
	d := &driver{
		pins:     make(map[string]*pin),
		channels: make(map[string]*channel),
		meta: hal.Metadata{
			Name:         "rpi",
			Description:  "hardware peripherals and GPIO channels on the base raspberry pi hardware",
			Capabilities: []hal.Capability{hal.Input, hal.Output, hal.PWM},
		},
	}
	for i := range validGPIOPins {
		p, err := factory(i)

		if err != nil {
			return nil, fmt.Errorf("can't build hal pin %d: %v", i, err)
		}
		name := fmt.Sprintf("GP%d", i)
		d.pins[name] = &pin{
			name:       name,
			number:     i,
			digitalPin: p,
		}
	}

	for _, p := range []int{0, 1} {
		ch := &channel{
			pin:       p,
			driver:    pd,
			frequency: s.PWMFreq * 100000,
			name:      fmt.Sprintf("%d", p),
		}
		d.channels[ch.name] = ch
	}
	return d, nil
}
