package displays

import (
	"github.com/gdamore/tcell"
)

var bits = []byte{0x80, 0x40, 0x20, 0x10}

type TextDisplay struct {
	screen tcell.Screen
	tCell  tcell.Style
	eCell  tcell.Style
	Memory [2048]bool
}

func (t *TextDisplay) Clear() {
	t.screen.Clear()
}

func (t *TextDisplay) Dispose() {
	t.screen.Fini()
}
func (t *TextDisplay) SetSprite(x uint8, y uint8, sprites []uint8) bool {
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

func (t *TextDisplay) GetPixel(x uint8, y uint8) bool {
	return t.Memory[int32(x)+int32(y)*64]
}

func (t *TextDisplay) SetPixel(x uint8, y uint8, on bool) bool {
	c := false

	if on && t.Memory[x+(y*64)] {
		c = true
	}
	x2 := int32(x) % 64
	t.Memory[x2+int32(y)*64] = on

	if t.GetPixel(x, y) {
		t.screen.SetCell(int(x2), int(y), t.tCell, ' ')
	} else {
		t.screen.SetCell(int(x2), int(y), t.eCell, ' ')
	}

	return c
}
func (t *TextDisplay) Render() {
	t.screen.Clear()
	for y := uint8(0); y < 32; y++ {
		for x := uint8(0); x < 64; x++ {

		}
	}

	t.screen.Sync()
}

// NewTextDisplay renders the screen using the library tcell.
// Still needs work
func NewTextDisplay() (*TextDisplay, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return &TextDisplay{}, err
	}
	err = screen.Init()
	if err != nil {
		return &TextDisplay{}, err

	}

	filledState := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite)
	emptyState := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack)

	screen.SetStyle(emptyState)

	return &TextDisplay{
		screen: screen,
		tCell:  filledState,
		eCell:  emptyState,
	}, nil
}
