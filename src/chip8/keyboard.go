package chip8

type TestKeyboard struct {
	KeysDown  []uint8
	WaitFor   uint8
	HasWaited bool
}

func (t *TestKeyboard) WaitForKey() uint8 {
	t.HasWaited = true
	return t.WaitFor
}

func (t *TestKeyboard) IsDown(key uint8) bool {
	for _, v := range t.KeysDown {
		if v == key {
			return true
		}
	}

	return false
}

type NoKeyboard struct {
}

func (n NoKeyboard) IsDown(_ uint8) bool {
	return false
}

func (n NoKeyboard) WaitForKey() uint8 {
	return 0
}
