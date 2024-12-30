package hal

type npin struct{}

func (n *npin) Close() error                 { return nil }
func (n *npin) SetDirection(_ bool) error    { return nil }
func (n *npin) Read() (int, error)           { return 0, nil }
func (n *npin) Write(_ int) error            { return nil }
func NoopPinFactory(int) (DigitalPin, error) { return new(npin), nil }
