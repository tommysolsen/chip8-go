package displays

import (
	"fmt"
)

type DebugDisplay struct {
	Memory [2048]bool
}

func (t *DebugDisplay) Clear() {
}

func (t *DebugDisplay) Dispose() {

}
func (t *DebugDisplay) SetSprite(x uint8, y uint8, sprites []uint8) bool {
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

func (t *DebugDisplay) GetPixel(x uint8, y uint8) bool {
	return t.Memory[int32(x)+int32(y)*64]
}

func (t *DebugDisplay) SetPixel(x uint8, y uint8, on bool) bool {
	fmt.Println(x, y, x+(y*uint8(64)), on)

	c := false

	if on && t.Memory[x+(y*64)] {
		c = true
	}
	t.Memory[(int32(x)%64)+(int32(y)*64)] = on

	return c
}

func (t *DebugDisplay) Render() {
	for y := uint8(0); y < 32; y++ {
		fmt.Printf("\r\n%02d:", y)
		for x := uint8(0); x < 64; x++ {
			//fmt.Printf("%02d,%02d ", x, y)

			if t.GetPixel(x, y) {
				fmt.Print("*")
			} else {
				fmt.Print("_")
			}
		}
	}

}

// NewDebugDisplay returns a new debug display which just crudely prints to screen
func NewDebugDisplay() (*DebugDisplay, error) {
	return &DebugDisplay{}, nil
}
