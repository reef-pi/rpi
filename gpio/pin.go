package gpio

type Pin struct {
	pin    uint8
	driver *Driver
}

// Set pin as Input
func (p *Pin) Input() {
	p.Direction(Input)
}

// Set pin as Output
func (p *Pin) Output() {
	p.Direction(Output)
}

// Set pin High
func (p *Pin) High() {
	p.Write(High)
}

// Set pin Low
func (p *Pin) Low() {
	p.Write(Low)
}

// Set a given pull up/down mode
func (p *Pin) Pull(pull Pull) {
	p.PullMode(pull)
}

// Pull up pin
func (p *Pin) PullUp() {
	p.PullMode(PullUp)
}

// Pull down pin
func (p *Pin) PullDown() {
	p.PullMode(PullDown)
}

// Disable pullup/down on pin
func (p *Pin) PullOff() {
	p.PullMode(PullOff)
}

// PinMode sets the direction of a given pin (Input or Output)
func (p *Pin) Direction(direction Direction) {
	p.driver.PinDirection(p.pin, direction)
}

// WritePin sets a given pin High or Low
// by setting the clear or set registers respectively
func (p *Pin) Write(state State) {
	p.driver.WriteToPin(p.pin, state)
}

// Read the state of a pin
func (p *Pin) Read() State {
	return p.driver.ReadFromPin(p.pin)
}

func (p *Pin) PullMode(pull Pull) {
	p.driver.PinPullMode(p.pin, pull)
}
