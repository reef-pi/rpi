package hal

import (
	"fmt"

	"github.com/reef-pi/rpi/i2c"

	"github.com/reef-pi/hal"
	pwmdriver "github.com/reef-pi/rpi/pwm"

	"github.com/kidoman/embd"
)

type Settings struct {
	PWMFreq int
}

type driver struct {
	pins     map[string]*rpiPin
	channels map[string]*rpiPwmChannel

	PinFactory func(key interface{}) (embd.DigitalPin, error)
	PWMFactory func() pwmdriver.Driver
}

func (r *driver) Metadata() hal.Metadata {
	return hal.Metadata{
		Name:        "rpi",
		Description: "hardware peripherals and GPIO channels on the base raspberry pi hardware",
		Capabilities: hal.Capabilities{
			Input:  true,
			Output: true,
			PWM:    true,
		},
	}
}

func (r *driver) Close() error {
	for _, pin := range r.pins {
		err := pin.Close()
		if err != nil {
			return fmt.Errorf("can't close hal driver due to channel %s", pin.Name())
		}
	}
	return nil
}

func (r *driver) init(s Settings) error {
	if r.PinFactory == nil {
		r.PinFactory = embd.NewDigitalPin
	}
	if r.PWMFactory == nil {
		r.PWMFactory = pwmdriver.New
	}
	if r.pins == nil {
		r.pins = make(map[string]*rpiPin)
	}
	if r.channels == nil {
		r.channels = make(map[string]*rpiPwmChannel)
	}

	for pin := range validGPIOPins {
		digitalPin, err := r.PinFactory(pin)

		if err != nil {
			return fmt.Errorf("can't build hal channel %d: %v", pin, err)
		}

		pin := rpiPin{
			name:       fmt.Sprintf("GP%d", pin),
			pin:        pin,
			digitalPin: digitalPin,
		}
		r.pins[pin.name] = &pin
	}

	pwmDriver := r.PWMFactory()

	for _, pin := range []int{0, 1} {
		pwmPin := &rpiPwmChannel{
			channel:   pin,
			driver:    pwmDriver,
			frequency: s.PWMFreq * 100000,
			name:      fmt.Sprintf("%d", pin),
		}
		r.channels[pwmPin.name] = pwmPin
	}

	return nil
}

func New(s Settings, b i2c.Bus) (hal.Driver, error) {
	d := &driver{}
	err := d.init(s)
	if err != nil {
		return nil, err
	}
	return d, nil
}
