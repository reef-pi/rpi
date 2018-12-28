package pwm

import (
	"os"
)

type recorder struct {
	values map[string][]byte
}

func (r *recorder) Get(s string) []byte {
	return r.values[s]
}

func Noop() (Driver, *recorder) {
	rec := &recorder{values: make(map[string][]byte)}
	d := &driver{
		sysfs: SysFS,
		writeFile: func(f string, c []byte, p os.FileMode) error {
			rec.values[f] = c
			return nil
		},
	}
	return d, rec
}
