package chip8

type NoDisplay struct {
}

func (n NoDisplay) Clear() {
}

func (n NoDisplay) SetSprite(_ uint8, _ uint8, _ []uint8) bool { return false }

func (n NoDisplay) GetPixel(_ uint8, _ uint8) bool {
	return false
}

func (n NoDisplay) SetPixel(_ uint8, _ uint8, _ bool) bool {
	return false
}

func (n NoDisplay) Render() {
}

type TestClearDisplay struct {
	*NoDisplay
	HasBeenCleared bool
}

func (t *TestClearDisplay) Clear() {
	t.HasBeenCleared = true
}
