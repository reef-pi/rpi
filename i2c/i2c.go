// +build !windows

package i2c

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"syscall"
	"unsafe"
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

type message struct {
	addr  uint16
	flags uint16
	len   uint16
	buf   uintptr
}

type ioctlData struct {
	msgs uintptr
	nmsg uint32
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
	f, err := os.OpenFile("/dev/i2c-1", os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return nil, err
	}
	return &bus{f: f, mu: new(sync.Mutex), syscallFn: syscall.Syscall}, nil
}

func (b *bus) send(cmd, addr uintptr) error {
	if _, _, errno := b.syscallFn(syscall.SYS_IOCTL, b.f.Fd(), cmd, addr); errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func (b *bus) Close() error {
	return b.f.Close()
}

func (b *bus) SetAddress(addr byte) error {
	return b.send(slaveCmd, uintptr(addr))
}

func (b *bus) ReadBytes(addr byte, num int) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := b.SetAddress(addr); err != nil {
		return []byte{0}, err
	}
	bytes := make([]byte, num)
	n, err := b.f.Read(bytes)
	if err != nil {
		return nil, err
	}
	if n != num {
		return []byte{0}, fmt.Errorf("i2c: Unexpected number (%v) of bytes read", n)
	}
	return bytes, nil
}

func (b *bus) WriteBytes(addr byte, value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := b.SetAddress(addr); err != nil {
		return err
	}
	_, err := b.f.Write(value)
	return err
}

func (b *bus) ReadFromReg(addr, reg byte, value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&value))

	var msgs [2]message
	msgs[0].addr = uint16(addr)
	msgs[0].flags = 0
	msgs[0].len = 1
	msgs[0].buf = uintptr(unsafe.Pointer(&reg))

	msgs[1].addr = uint16(addr)
	msgs[1].flags = rd
	msgs[1].len = uint16(len(value))
	msgs[1].buf = uintptr(unsafe.Pointer(hdrp.Data))

	var d ioctlData

	d.msgs = uintptr(unsafe.Pointer(&msgs))
	d.nmsg = 2
	if err := b.SetAddress(addr); err != nil {
		return err
	}
	return b.send(rdrwCmd, uintptr(unsafe.Pointer(&d)))
}

func (b *bus) WriteToReg(addr, reg byte, value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	outbuf := append([]byte{reg}, value...)

	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&outbuf))

	var msg message
	msg.addr = uint16(addr)
	msg.flags = 0
	msg.len = uint16(len(outbuf))
	msg.buf = uintptr(unsafe.Pointer(hdrp.Data))

	var d ioctlData
	d.msgs = uintptr(unsafe.Pointer(&msg))
	d.nmsg = 1

	if err := b.SetAddress(addr); err != nil {
		return err
	}
	return b.send(rdrwCmd, uintptr(unsafe.Pointer(&d)))
}
