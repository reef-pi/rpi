// +build !windows

package hal

import "github.com/kidoman/embd"

func newDigitalPin(int i) (DigitalPin, err) {
	return embd.NewDigitalPin(i)
}
