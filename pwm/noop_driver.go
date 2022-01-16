package pwm

import (
	"os"
	"sync"
)

type recorder struct {
	mu     sync.Mutex
	values map[string][]byte
}

func (r *recorder) Get(s string) []byte {
	return r.values[s]
}

func Noop() (Driver, *recorder) {
	rec := &recorder{values: make(map[string][]byte), mu: sync.Mutex{}}
	d := &driver{
		sysfs: SysFS,
		writeFile: func(f string, c []byte, p os.FileMode) error {
			rec.mu.Lock()
			defer rec.mu.Unlock()
			rec.values[f] = c
			return nil
		},
		readFile: func(f string) ([]byte, error) {
			rec.mu.Lock()
			_, ok := rec.values[f]
			rec.mu.Unlock()
			if !ok {
				return []byte{}, nil
			}
			return rec.values[f], nil
		},
	}
	return d, rec
}
