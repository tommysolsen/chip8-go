package displays

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

var binds = map[sdl.Scancode]uint8{
	sdl.SCANCODE_1: 1,
	sdl.SCANCODE_2: 2,
	sdl.SCANCODE_3: 3,
	sdl.SCANCODE_Q: 4,
	sdl.SCANCODE_W: 5,
	sdl.SCANCODE_E: 6,
	sdl.SCANCODE_A: 7,
	sdl.SCANCODE_S: 8,
	sdl.SCANCODE_D: 9,
	sdl.SCANCODE_Z: 0x0A,
	sdl.SCANCODE_X: 0,
	sdl.SCANCODE_C: 0x0B,
	sdl.SCANCODE_4: 0x0C,
	sdl.SCANCODE_R: 0x0D,
	sdl.SCANCODE_F: 0x0E,
	sdl.SCANCODE_V: 0x0F,
}

type NewSDLDisplay struct {
	window    *sdl.Window
	renderer  *sdl.Renderer
	pixelSize int32

	Memory       [4096]bool
	renderNeeded bool
}

func (t *NewSDLDisplay) IsDown(key uint8) bool {
	keys := sdl.GetKeyboardState()

	for _, i := range keys {
		k := sdl.GetScancodeFromKey(sdl.Keycode(i))
		if binds[k] == key {
			return true
		}
	}

	return false
}

func (t *NewSDLDisplay) WaitForKey() uint8 {
	for true {
		fmt.Println("Waiting for key")
		switch s := sdl.PollEvent().(type) {
		case *sdl.KeyboardEvent:
			if s.State != sdl.PRESSED {
				break
			}
			for key, v := range binds {
				if s.Keysym.Scancode == key {
					return v
				}
			}
			break
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	panic("Should never happen")
}

func (t *NewSDLDisplay) Clear() {
	rect := sdl.Rect{
		X: 0,
		Y: 0,
		W: 64 * t.pixelSize,
		H: 32 * t.pixelSize,
	}
	_ = t.renderer.FillRect(&rect)
	t.renderer.Present()
}

func (t *NewSDLDisplay) Dispose() {
	_ = t.window.Destroy()
	_ = t.renderer.Destroy()
}

func (t *NewSDLDisplay) SetSprite(x uint8, y uint8, sprites []uint8) bool {
	c := false
	for i, sprite := range sprites {
		for i2, bit := range bits {
			if t.SetPixel(x+uint8(i2), y+uint8(i), sprite&bit == bit) {
				c = true
			}
		}
	}

	return c
}

func (t *NewSDLDisplay) GetPixel(x uint8, y uint8) bool {
	return t.Memory[memoryLocation(x, y)]
}

func (t *NewSDLDisplay) SetPixel(x uint8, y uint8, on bool) bool {
	mLoc := memoryLocation(x, y)
	c := false

	if on && t.Memory[mLoc] {
		c = true
	}

	t.renderNeeded = t.renderNeeded || on != t.Memory[mLoc]
	if on {
		t.Memory[mLoc] = !t.Memory[mLoc]
	}

	return c
}

func (t *NewSDLDisplay) Render() {
	if !t.renderNeeded {
		return
	}
	_ = t.renderer.Clear()
	for x := 0; x < 64; x++ {
		for y := 0; y < 32; y++ {
			rect := sdl.Rect{
				X: int32(x) % 64 * t.pixelSize,
				Y: int32(y) * t.pixelSize,
				W: t.pixelSize,
				H: t.pixelSize,
			}

			if t.GetPixel(uint8(x), uint8(y)) {
				_ = t.renderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
			} else {
				_ = t.renderer.SetDrawColor(0x00, 0x00, 0x00, 0xFF)
			}
			_ = t.renderer.FillRect(&rect)

		}
	}
	t.renderer.Present()
	t.renderNeeded = false
}

func NewSDLRenderer(pixelSize int32) (*NewSDLDisplay, error) {
	window, err := sdl.CreateWindow("Chip8", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 64*pixelSize, 32*pixelSize, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}
	return &NewSDLDisplay{
		pixelSize: pixelSize,
		window:    window,
		renderer:  renderer,
	}, nil
}
