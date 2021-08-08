package chip8

type StateDumper interface {
	DumpState(c Cpu)
}

type RngGenerator interface {
	GetRandom() uint8
}

type Display interface {
	GetPixel(x uint8, y uint8) bool
	SetPixel(x uint8, y uint8, on bool) bool
	SetSprite(x uint8, y uint8, sprite []uint8) bool
	Clear()
	Render()
}

type Keyboard interface {
	IsDown(key uint8) bool
	WaitForKey() uint8
}
