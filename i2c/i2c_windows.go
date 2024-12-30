//go:build windows
// +build windows

package i2c

import (
	"io"
	"sync"
	"syscall"
)

const (
	delay    = 20
	slaveCmd = 0x0703 // Cmd to set slave address
	rdrwCmd  = 0x0707 // Cmd to read/write data together
	rd       = 0x0001
)

type Fd interface {
	io.ReadWriteCloser
	Fd() uintptr
}

type Bus interface {
	SetAddress(addr byte) error
	ReadBytes(addr byte, num int) ([]byte, error)
	WriteBytes(addr byte, value []byte) error
	ReadFromReg(addr, reg byte, value []byte) error
	WriteToReg(addr, reg byte, value []byte) error
	Close() error
}

type bus struct {
	f         Fd
	syscallFn func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
	mu        *sync.Mutex
}

func New() (*bus, error) {
	return &bus{}, nil
}

func (b *bus) Close() error {
	return nil
}

func (b *bus) SetAddress(addr byte) error {
	return nil
}

func (b *bus) ReadBytes(addr byte, num int) ([]byte, error) {
	return nil, nil
}

func (b *bus) WriteBytes(addr byte, value []byte) error {
	return nil
}

func (b *bus) ReadFromReg(addr, reg byte, value []byte) error {
	return nil
}

func (b *bus) WriteToReg(addr, reg byte, value []byte) error {
	return nil
}
