package hal

import (
	"fmt"
	"sort"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/pwm"
)

type channel struct {
	pin       int
	name      string
	driver    pwm.Driver
	frequency int
}

func (p *channel) Set(value float64) error {
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

	setting := float64(p.frequency/1000) * value
	if err := p.driver.DutyCycle(p.pin, int(setting)); err != nil {
		return err
	}
	return p.driver.Enable(p.pin)
}

func (p *channel) Name() string {
	return p.name
}

func (r *driver) PWMChannels() []hal.PWMChannel {
	var chs []hal.PWMChannel
	for _, ch := range r.channels {
		chs = append(chs, ch)
	}
	sort.Slice(chs, func(i, j int) bool { return chs[i].Name() < chs[j].Name() })
	return chs
}

func (r *driver) PWMChannel(name string) (hal.PWMChannel, error) {
	ch, ok := r.channels[name]
	if !ok {
		return nil, fmt.Errorf("unknown pwm channel %s", name)
	}
	return ch, nil
}
