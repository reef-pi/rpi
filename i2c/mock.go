package i2c

type mock struct{}

func (m *mock) SetAddress(_ byte) error { return nil }
func (m *mock) ReadBytes(addr byte, num int) ([]byte, error) {
	return make([]byte, num, num), nil
}
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
