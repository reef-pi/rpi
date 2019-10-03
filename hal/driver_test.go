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
	for _, expected := range []hal.Capability{hal.DigitalInput, hal.DigitalOutput, hal.PWM} {
		if !meta.HasCapability(expected) {
			t.Error("didn't find expected capabilities")
		}
	}
	for _, cap := range meta.Capabilities {
		if cap == hal.AnalogInput {
			t.Error("rpi can't provide pH")
		}
	}

	input := hal.DigitalInputDriver(d)

	pins := input.DigitalInputPins()
	if l := len(validGPIOPins); l != len(pins) {
		t.Error("Wrong pin count. Expected:", len(validGPIOPins), " found:", len(d.pins))
	}

	var output hal.DigitalOutputDriver = d
	outPins := output.DigitalOutputPins()
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

	input := hal.DigitalInputDriver(d)
	output := hal.DigitalOutputDriver(d)

	ipins := input.DigitalInputPins()
	opins := output.DigitalOutputPins()
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
	output := hal.DigitalOutputDriver(d)
	pin, err := output.DigitalOutputPin(26)
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

	ch, err := pwmDriver.PWMChannel(0)
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

func TestPinMap(t *testing.T) {
	s := Settings{}
	s.PWMFreq = 100
	pd, _ := pwm.Noop()
	d, err := NewAdapter(s, pd, NoopPinFactory)
	if err != nil {
		t.Error(err)
	}
	iPins, err := d.Pins(hal.DigitalInput)
	if err != nil {
		t.Error(err)
	}
	if len(iPins) != 26 {
		t.Error("Expected 26 digital input pins. Found:", len(iPins))
	}
	oPins, err := d.Pins(hal.DigitalOutput)
	if err != nil {
		t.Error(err)
	}
	if len(oPins) != 26 {
		t.Error("Expected 26 digital output pins. Found:", len(oPins))
	}
	pPins, err := d.Pins(hal.PWM)
	if err != nil {
		t.Error(err)
	}
	if len(pPins) != 2 {
		t.Error("Expected 2 pwm pins. Found:", len(pPins))
	}
}
