package hal

import (
	"fmt"
	"log"
	"sort"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/pwm"
)

type channel struct {
	pin       int
	name      string
	driver    pwm.Driver
	frequency int
	v         float64
}

func (p *channel) Set(value float64) error {
	if p.frequency <= 0 {
		log.Printf("warning: RPI PWM frequency is 0, defaulting to 150")
		p.frequency = 150
	}
	if value < 0 || value > 100 {
		return fmt.Errorf("value must be 0-100, got %f", value)
	}

	exported, err := p.driver.IsExported(p.pin)
	if err != nil {
		return err
	}
	if !exported {
		if err := p.driver.Export(p.pin); err != nil {
			return err
		}
	}
	if err := p.driver.Frequency(p.pin, p.frequency); err != nil {
		return err
	}
	if err := p.driver.DutyCycle(p.pin, value); err != nil {
		return err
	}
	if err := p.driver.Enable(p.pin); err != nil {
		return err
	}
	p.v = value
	return nil
}

func (ch *channel) Close() error { return nil }
func (ch *channel) LastState() bool {
	return ch.v == 100
}

func (ch *channel) Write(b bool) error {
	var v float64
	if b == true {
		v = 100
	}
	return ch.Set(v)
}

func (p *channel) Name() string {
	return p.name
}

func (p *channel) Number() int {
	return p.pin
}
func (r *driver) PWMChannels() []hal.PWMChannel {
	var chs []hal.PWMChannel
	for _, ch := range r.channels {
		chs = append(chs, ch)
	}
	sort.Slice(chs, func(i, j int) bool { return chs[i].Name() < chs[j].Name() })
	return chs
}

func (r *driver) PWMChannel(p int) (hal.PWMChannel, error) {
	ch, ok := r.channels[p]
	if !ok {
		return nil, fmt.Errorf("unknown pwm channel %d", p)
	}
	return ch, nil
}
