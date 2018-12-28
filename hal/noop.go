package hal

import (
	"github.com/kidoman/embd"
)

type npin struct{}

func (n *npin) Close() error                           { return nil }
func (n *npin) SetDirection(_ embd.Direction) error    { return nil }
func (n *npin) Read() (int, error)                     { return 0, nil }
func (n *npin) Write(_ int) error                      { return nil }
func NoopPinFactory(_ interface{}) (DigitalPin, error) { return new(npin), nil }
