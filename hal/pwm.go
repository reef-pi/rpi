package hal

import (
	"fmt"
	"sort"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/pwm"
)

type rpiPwmChannel struct {
	channel   int
	name      string
	driver    pwm.Driver
	frequency int
}

func (p *rpiPwmChannel) Set(value float64) error {
	if value < 0 || value > 100 {
		return fmt.Errorf("value must be 0-100, got %f", value)
	}

	exported, err := p.driver.IsExported(p.channel)
	if err != nil {
		return err
	}
	if !exported {
		if err := p.driver.Export(p.channel); err != nil {
			return err
		}
	}
	if err := p.driver.Frequency(p.channel, p.frequency); err != nil {
		return err
	}

	setting := float64(p.frequency/1000) * value
	if err := p.driver.DutyCycle(p.channel, int(setting)); err != nil {
		return err
	}
	return p.driver.Enable(p.channel)
}

func (p *rpiPwmChannel) Name() string {
	return p.name
}

func (r *driver) Channels() []hal.Channel {
	var chs []hal.Channel
	for _, ch := range r.channels {
		chs = append(chs, ch)
	}
	sort.Slice(chs, func(i, j int) bool { return chs[i].Name() < chs[j].Name() })
	return chs
}

func (r *driver) GetChannel(name string) (hal.Channel, error) {
	ch, ok := r.channels[name]
	if !ok {
		return nil, fmt.Errorf("unknown pwm channel %s", name)
	}
	return ch, nil
}
