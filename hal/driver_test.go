package hal

import (
	"testing"

	"github.com/reef-pi/hal"
)

func TestNewRPiDriver(t *testing.T) {
	s := Settings{}
	s.PWMFreq = 100
	d, err := New(s, newMockPWMDriver(), newMockDigitalPin)
	if err != nil {
		t.Error(err)
	}

	meta := d.Metadata()
	if meta.Name != "rpi" {
		t.Error("driver name wasn't rpi")
	}
	if !(meta.Capabilities.PWM &&
		meta.Capabilities.Input &&
		meta.Capabilities.Output) {
		t.Error("didn't find expected capabilities")
	}
	if meta.Capabilities.PH {
		t.Error("rpi can't provide pH")
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
	for _, i := range d.pins {
		p := i.digitalPin.(*mockDigitalPin)
		if !p.closed {
			t.Errorf("pin %v wasn't closed", p)
		}
	}
}

func TestRpiDriver_InputPins(t *testing.T) {
	s := Settings{}
	s.PWMFreq = 100
	d, err := New(s, newMockPWMDriver(), newMockDigitalPin)
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

	v, err = ipins[1].Read()
	if err != nil {
		t.Errorf("unexpected error reading pin %v", err)
	}
	if !v {
		t.Error("expected pin to be high")
	}
}

func TestRpiDriver_GetOutputPin(t *testing.T) {
	s := Settings{}
	s.PWMFreq = 100
	d, err := New(s, newMockPWMDriver(), newMockDigitalPin)
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
	d, err := New(s, newMockPWMDriver(), newMockDigitalPin)
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

	backingChannel := ch.(*channel)
	backingDriver := backingChannel.driver.(*mockPwmDriver)

	if s := backingDriver.setting[0]; s != 100000 {
		t.Errorf("backing driver not reporting 100000, got %d", s)
	}
	if !backingDriver.enabled[0] {
		t.Error("backing driver was not enabled")
	}

}
