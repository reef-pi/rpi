package pwm

// http://www.jumpnowtek.com/rpi/Using-the-Raspberry-Pi-Hardware-PWM-timers.html

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	}
}

func toS(ch int) []byte {
	return []byte(fmt.Sprintf("%d\n", ch))
}

type driver struct {
	writeFile func(file string, data []byte, perm os.FileMode) error
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

func (d *driver) DutyCycle(ch, duty int) error {
	file := filepath.Join(d.sysfs, fmt.Sprintf("pwm%d", ch), "duty_cycle")
	return d.writeFile(file, toS(duty), 0644)
}

func (d *driver) Frequency(ch, freq int) error {
	file := filepath.Join(d.sysfs, fmt.Sprintf("pwm%d", ch), "period")
	return d.writeFile(file, toS(freq), 0644)
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
