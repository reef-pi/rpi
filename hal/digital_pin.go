// +build !windows

package hal

import "github.com/reef-pi/embd"

func newDigitalPin(i int) (DigitalPin, error) {
	return embd.NewDigitalPin(i)
}
