package gpio

import (
	"reflect"
	"sync"
	"time"
	"unsafe"
)

type Direction uint8

const (
	Input Direction = iota
	Output
)

type State uint8

const (
	Low State = iota
	High
)

type Pull uint8

const (
	PullOff Pull = iota
	PullDown
	PullUp
)

type Driver struct {
	mu   sync.Mutex
	mem  []uint32
	mem8 []uint8
}

// Create from memory mapped data
func CreateFromMmap(mem8 []uint8) *Driver {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&mem8))
	header.Len /= (32 / 8) // (32 bit = 4 bytes)
	header.Cap /= (32 / 8)
	mem := *(*[]uint32)(unsafe.Pointer(&header))
	driver := Driver{
		mu:   sync.Mutex{},
		mem8: mem8,
		mem:  mem,
	}
	return &driver
}

// Close unmaps GPIO memory
func (d *Driver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return munmap(d.mem8)
}

func (d *Driver) Pin(pin uint8) *Pin {
	return &Pin{pin: pin, driver: d}
}

func (d *Driver) PinDirection(pin uint8, direction Direction) {
	pinMask := uint32(7) // 0b111 - pinmode is 3 bits
	// Pin fsel register, 0 or 1 depending on bank
	fsel := uint8(pin) / 10
	shift := (uint8(pin) % 10) * 3

	v := (d.mem[fsel] &^ (pinMask << shift)) | (1 << shift)
	if direction == Input {
		d.mem[fsel] = d.mem[fsel] &^ (pinMask << shift)
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.mem[fsel] = v
}

func (d *Driver) WriteToPin(pin uint8, state State) {
	p := uint8(pin)
	// Set register, 7 / 8 depending in high state
	// Clear register, 10 / 11 depending in low state
	reg := p/32 + 7
	if state == Low {
		reg = p/32 + 10
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.mem[reg] = 1 << (p & 31)
}

func (d *Driver) ReadFromPin(pin uint8) State {
	// Input level register offset (13 / 14 depending on state)
	reg := pin/32 + 13
	if (d.mem[reg] & (1 << pin)) != 0 {
		return High
	}
	return Low
}

func (d *Driver) PinPullMode(pin uint8, pull Pull) {
	// Pull up/down/off register has offset 38 / 39, pull is 37
	pullClkReg := pin/32 + 38
	pullReg := 37
	shift := pin % 32
	d.mu.Lock()
	defer d.mu.Unlock()
	switch pull {
	case PullDown, PullUp:
		d.mem[pullReg] = d.mem[pullReg]&^3 | uint32(pull)
	case PullOff:
		d.mem[pullReg] = d.mem[pullReg] &^ 3
	}
	// Wait for value to clock in, this is ugly, sorry :(
	time.Sleep(time.Microsecond)
	d.mem[pullClkReg] = 1 << shift
	// Wait for value to clock in
	time.Sleep(time.Microsecond)
	d.mem[pullReg] = d.mem[pullReg] &^ 3
	d.mem[pullClkReg] = 0
}
