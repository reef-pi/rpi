//go:build windows
// +build windows

package hal

import (
	"time"
)

func newDigitalPin(i int) (DigitalPin, error) {
	pin := mockPin{i}
	return pin, nil
}

type mockPin struct {
	pinNumber int
}

func (p mockPin) InterruptPin()                              {}
func (p mockPin) N() int                                     { return p.pinNumber }
func (p mockPin) Write(_ int) error                          { return nil }
func (p mockPin) Read() (int, error)                         { return 0, nil }
func (p mockPin) TimePulse(state int) (time.Duration, error) { return time.Duration(0), nil }
func (p mockPin) SetDirection(dir string) error              { return nil }
func (p mockPin) ActiveLow(b bool) error                     { return nil }
func (p mockPin) PullUp() error                              { return nil }
func (p mockPin) PullDown() error                            { return nil }
func (p mockPin) Close() error                               { return nil }
