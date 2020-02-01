package pwm

// http://www.jumpnowtek.com/rpi/Using-the-Raspberry-Pi-Hardware-PWM-timers.html

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	SysFS = "/sys/class/pwm/pwmchip0"
)

type Driver interface {
	Export(ch int) error
	Unexport(ch int) error
	DutyCycle(ch, duty int) error
	Frequency(ch, freq int) error
	Enable(ch int) error
	Disable(ch int) error
	IsEnabled(ch int) (bool, error)
	IsExported(ch int) (bool, error)
}

func New() Driver {
	return &driver{
		sysfs:     SysFS,
		writeFile: ioutil.WriteFile,
		readFile:  ioutil.ReadFile,
	}
}

func toS(ch int) []byte {
	return []byte(fmt.Sprintf("%d\n", ch))
}

func toS64(value int64) []byte {
	return []byte(fmt.Sprintf("%d\n", value))
}

type driver struct {
	writeFile func(file string, data []byte, perm os.FileMode) error
	readFile  func(file string) ([]byte, error)
	sysfs     string
}

func (d *driver) Export(ch int) error {
	file := filepath.Join(d.sysfs, "export")
	return d.writeFile(file, toS(ch), 0600)
}

func (d *driver) Unexport(ch int) error {
	file := filepath.Join(d.sysfs, "unexport")
	return d.writeFile(file, toS(ch), 0600)
}

// DutyCycle sets the duty cycle as a 0-100 percentage of the
// given period of the driver. Does not assume the Frequency
// method has been called.
func (d *driver) DutyCycle(ch, duty int) error {
	freqFile := filepath.Join(d.sysfs, fmt.Sprintf("pwm%d", ch), "period")
	data, err := d.readFile(freqFile)
	if err != nil {
		return err
	}
	sdata := strings.TrimRight(string(data), "\n")
	period, err := strconv.ParseInt(sdata, 10, 64)
	if err != nil {
		return err
	}
	nanoSecondsDuty := int64((float64(duty) / 100.0) * float64(period))

	file := filepath.Join(d.sysfs, fmt.Sprintf("pwm%d", ch), "duty_cycle")
	return d.writeFile(file, toS64(nanoSecondsDuty), 0644)
}

// Frequency sets the frequency in Hz. This is written as nano-seconds to the
// period register in the underlying sysfs pwm driver.
// Note: the kernel driver will refuse to write period > duty_cycle
// so we need to check this first, and reset the duty_cycle to 0 if this
// is the case.
func (d *driver) Frequency(ch, freq int) error {
	period := int64((1.0 / (float64(freq))) * 1.0e9)
	dutyCycleFile := filepath.Join(d.sysfs, fmt.Sprintf("pwm%d", ch), "duty_cycle")
	data, err := d.readFile(dutyCycleFile)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		sdata := strings.TrimRight(string(data), "\n")
		dutyCycle, err := strconv.ParseInt(sdata, 10, 64)
		if err == nil {
			if dutyCycle > period {
				err = d.writeFile(dutyCycleFile, toS64(0), 0644)
				if err != nil {
					return err
				}
			}
		}
	}
	file := filepath.Join(d.sysfs, fmt.Sprintf("pwm%d", ch), "period")
	return d.writeFile(file, toS64(period), 0644)
}

func (d *driver) Enable(ch int) error {
	file := filepath.Join(d.sysfs, fmt.Sprintf("pwm%d", ch), "enable")
	return d.writeFile(file, toS(1), 0644)
}

func (d *driver) Disable(ch int) error {
	file := filepath.Join(d.sysfs, fmt.Sprintf("pwm%d", ch), "enable")
	return d.writeFile(file, toS(0), 0644)
}
func (d *driver) IsEnabled(ch int) (bool, error) {
	file := filepath.Join(d.sysfs, fmt.Sprintf("pwm%d", ch), "enable")
	v, err := ioutil.ReadFile(file)
	if err != nil {
		return false, err
	}
	s := strings.TrimSpace(string(v))
	return s == "1", nil
}
func (d *driver) IsExported(ch int) (bool, error) {
	file := filepath.Join(d.sysfs, fmt.Sprintf("pwm%d", ch), "enable")
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
