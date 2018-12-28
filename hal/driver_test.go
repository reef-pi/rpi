package hal

import (
	"fmt"
	"testing"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/pwm"
	"path/filepath"
)

var s = Settings{
	PWMFreq: 100,
}

func mockPWMDriver() pwm.Driver {
	d, _ := pwm.Noop()
	return d
}

func TestNewRPiDriver(t *testing.T) {
	d, err := NewAdapter(s, mockPWMDriver(), NoopPinFactory)
	if err != nil {
		t.Error(err)
	}

	meta := d.Metadata()
	if meta.Name != "rpi" {
		t.Error("driver name wasn't rpi")
	}
	for _, expected := range []hal.Capability{hal.Input, hal.Output, hal.PWM} {
		if !meta.HasCapability(expected) {
			t.Error("didn't find expected capabilities")
		}
	}
	for _, cap := range meta.Capabilities {
		if cap == hal.PH {
			t.Error("rpi can't provide pH")
		}
	}

	input := hal.InputDriver(d)

	pins := input.InputPins()
	if l := len(validGPIOPins); l != len(pins) {
		t.Error("Wrong pin count. Expected:", len(validGPIOPins), " found:", len(d.pins))
	}

	var output hal.OutputDriver = d
	outPins := output.OutputPins()
	if l := len(validGPIOPins); l != len(outPins) {
		t.Error("Wrong pin count. Expected:", len(validGPIOPins), " found:", len(outPins))
	}

	if err := d.Close(); err != nil {
		t.Errorf("unexpected error closing driver %v", err)
	}
}

func TestRpiDriver_InputPins(t *testing.T) {
	d, err := NewAdapter(s, mockPWMDriver(), NoopPinFactory)
	if err != nil {
		t.Error(err)
	}

	input := hal.InputDriver(d)
	output := hal.OutputDriver(d)

	ipins := input.InputPins()
	opins := output.OutputPins()
	if ipins[0].Name() != opins[0].Name() {
		t.Error("input and output pins don't match")
	}

	v, err := ipins[0].Read()
	if err != nil {
		t.Errorf("unexpected error reading pin %v", err)
	}
	if v {
		t.Error("expected pin to be low")
	}
	err = opins[1].Write(true)
	if err != nil {
		t.Errorf("unexpected error writing pin %v", err)
	}

	if _, err = ipins[1].Read(); err != nil {
		t.Errorf("unexpected error reading pin %v", err)
	}
}

func TestRpiDriver_GetOutputPin(t *testing.T) {
	d, err := NewAdapter(s, mockPWMDriver(), NoopPinFactory)
	if err != nil {
		t.Error(err)
	}
	output := hal.OutputDriver(d)
	pin, err := output.OutputPin("GP26")
	if err != nil {
		t.Errorf("could not get output pin %v", err)
	}
	if pin.Name() != "GP26" {
		t.Errorf("pin name %s was not GP26", pin.Name())
	}
}

func TestRpiDriver_GetPWMChannel(t *testing.T) {
	s := Settings{}
	s.PWMFreq = 100
	pd, rec := pwm.Noop()
	d, err := NewAdapter(s, pd, NoopPinFactory)
	if err != nil {
		t.Error(err)
	}
	pwmDriver := hal.PWMDriver(d)

	ch, err := pwmDriver.PWMChannel("0")
	if err != nil {
		t.Errorf("unexpected error getting pwm channel %v", err)
	}
	if name := ch.Name(); name != "0" {
		t.Error("PWM channel was not named '0'")
	}

	err = ch.Set(10)
	if err != nil {
		t.Errorf("unexpected error setting PWM %v", err)
	}

	file := filepath.Join(pwm.SysFS, "pwm0", "period")
	f := 10000000
	if s := rec.Get(file); string(s) != fmt.Sprintf("%d\n", f) {
		t.Errorf("backing driver not reporting %d, got %s", f, string(s))
	}

}
