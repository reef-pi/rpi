package gpio

import (
	"testing"
)

func TestPinInput(t *testing.T) {
	driver := CreateFromMmap(make([]uint8, 200))
	pin := driver.Pin(4)
	pin.Input()
	pin.Output()
	pin.Low()
	pin.PullUp()
	pin.PullDown()
	pin.PullOff()
	pin.High()
	pin.Low()
}
