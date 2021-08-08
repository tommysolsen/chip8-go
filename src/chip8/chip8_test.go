package chip8

import (
	"strings"
	"testing"
)

func TestCpu_LoadProgram(t *testing.T) {

	testString := "01234"

	cpu := NewCPU(0x200+int16(len(testString)), NoDisplay{}, nil)
	err := cpu.LoadProgram(strings.NewReader(testString))

	if err != nil {
		t.Errorf("Failed to load program: %s", err)
	}

	for i, c := range testString {
		if cpu.Memory[0x200+i] != uint8(c) {
			t.Errorf("Memory address 0x%x has invalid data: Wanted %d, got %d", 0x200+i, c, cpu.Memory[0x200+i])
		}
	}

}

func TestCpu_decrementTimers(t *testing.T) {
	cpu := NewCPU(4000, NoDisplay{}, nil)
	cpu.DT = 40
	cpu.ST = 39

	cpu.DecrementTimers()

	if cpu.DT != 39 {
		t.Errorf("Expected DT to be 39, found %d", cpu.DT)
	}

	if cpu.ST != 38 {
		t.Errorf("Expected DT to be 38, found %d", cpu.DT)
	}

	cpu.DT = 0
	cpu.DecrementTimers()
	if cpu.DT != 0 {
		t.Errorf("Expected DT to be 0, found %d", cpu.DT)
	}

	if cpu.ST != 37 {
		t.Errorf("Expected DT to be 37, found %d", cpu.DT)
	}
}
