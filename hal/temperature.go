package hal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/reef-pi/hal"
)

type temperatureChannel struct {
	path string
}

func (t *temperatureChannel) ReadTemperature() (float64, error) {
	fi, err := os.Open(filepath.Join(t.path, "w1_slave"))
	if err != nil {
		return -1, err
	}
	defer fi.Close()

	reader := bufio.NewReader(fi)
	l1, _, err := reader.ReadLine()
	if err != nil {
		return -1, err
	}
	if !strings.HasSuffix(string(l1), "YES") {
		return -1, fmt.Errorf("First line of device file does not ends with YES")
	}
	l2, _, err := reader.ReadLine()
	if err != nil {
		return -1, err
	}
	vals := strings.Split(string(l2), "=")
	if len(vals) < 2 {
		return -1, fmt.Errorf("Second line of device file does not have '=' separated temperature value")
	}
	v, err := strconv.Atoi(vals[1])
	if err != nil {
		return -1, err
	}
	temp := float64(v) / 1000.0
	return temp, nil
}

func (d *driver) TemperatureChannels() []hal.TemperatureChannel {
	files, err := filepath.Glob("/sys/bus/w1/devices/28-*")
	if err != nil {
		return nil
	}
	sort.Strings(files)
	var channels []hal.TemperatureChannel
	for _, p := range files {
		channels = append(channels, &temperatureChannel{p})
	}
	return channels
}

func (d *driver) TemperatureChannel(ch int) (hal.TemperatureChannel, error) {
	channels := d.TemperatureChannels()
	if ch > len(channels) {
		return nil, fmt.Errorf("invalid temperature sensor %d", ch)
	}
	return channels[ch-1], nil
}