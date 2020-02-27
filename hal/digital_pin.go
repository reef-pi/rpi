// +build !windows

package hal

import "github.com/kidoman/embd"

func newDigitalPin(i int) (DigitalPin, error) {
	return embd.NewDigitalPin(i)
}
