package i2c

import (
	"sync"
	"syscall"
	"testing"
)

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
