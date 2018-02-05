package i2c

import (
	"sync"
	"syscall"
	"testing"
)

type mock struct{}

func (m *mock) SetAddress(_ byte) error                        { return nil }
func (m *mock) ReadBytes(addr byte, num int) ([]byte, error)   { return []byte{}, nil }
func (m *mock) WriteBytes(addr byte, value []byte) error       { return nil }
func (m *mock) ReadFromReg(addr, reg byte, value []byte) error { return nil }
func (m *mock) WriteToReg(addr, reg byte, value []byte) error  { return nil }

func MockBus() Bus { return new(mock) }

type mockFs struct{}

func (m *mockFs) Fd() uintptr {
	return 1
}
func (m *mockFs) Read(b []byte) (int, error) {
	return len(b), nil
}
func (m *mockFs) Write(b []byte) (int, error) {
	return len(b), nil
}

func (m *mockFs) Close() error {
	return nil
}

func TestI2c(t *testing.T) {
	mockSycall := func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
		return 0, 0, 0
	}
	b := &bus{
		f:         &mockFs{},
		mu:        &sync.Mutex{},
		syscallFn: mockSycall,
	}
	if err := b.WriteBytes(byte(1), []byte("")); err != nil {
		t.Error(err)
	}
}
