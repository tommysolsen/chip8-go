package chip8

import (
	"bytes"
	"fmt"
	"io"
)

type Cpu struct {
	Memory []uint8
	V      [0x10]uint8
	PC     uint16
	SP     uint16
	S      [16]uint16
	I      uint16

	DT uint8
	ST uint8

	rng      RngGenerator
	display  Display
	keyboard Keyboard
}

func (cpu *Cpu) LoadProgram(program io.Reader) error {
	err := cpu.LoadCode(program, 0x200)
	if err != nil {
		return fmt.Errorf("Unable to load program: %s", err)
	}

	return nil
}

func (cpu *Cpu) DecrementTimers() {
	if cpu.ST > 0 {
		cpu.ST--
	}

	if cpu.DT > 0 {
		cpu.DT--
	}
}

func (cpu *Cpu) LoadCode(program io.Reader, from uint16) error {
	buffer := make([]byte, 100)

	var offset uint16 = 0
	for {
		n, err := program.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if n > 0 {
			for _, v := range buffer[0:n] {
				cpu.Memory[from+offset] = v
				offset += 1
			}
		}
	}

	fmt.Printf("Loaded %d bytes\r\n", offset)

	return nil
}

func (cpu *Cpu) LoadInterpreter() {
	_ = cpu.LoadCode(bytes.NewReader([]byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0,

		0x20, 0x60, 0x20, 0x20, 0x70,
		0xF0, 0x10, 0xF0, 0x80, 0xF0,
		0xF0, 0x10, 0xF0, 0x10, 0xF0,

		0x90, 0x90, 0xF0, 0x10, 0x10,
		0xF0, 0x10, 0x20, 0x40, 0x40,
		0xF0, 0x80, 0xF0, 0x90, 0xF0,

		0xF0, 0x10, 0x20, 0x40, 0x40,
		0xF0, 0x90, 0xF0, 0x90, 0xF0,
		0xF0, 0x90, 0xF0, 0x10, 0xF0,

		0xF0, 0x90, 0xF0, 0x90, 0x90,
		0xE0, 0x90, 0xE0, 0x90, 0xE0,
		0xF0, 0x80, 0x80, 0x80, 0xF0,

		0xE0, 0x90, 0x90, 0x90, 0xE0,
		0xF0, 0x80, 0xF0, 0x80, 0xF0,
		0xF0, 0x80, 0xF0, 0x80, 0x80,
	}), 0x0)

}

func NewCPU(memorySize int16, display Display, keyboard Keyboard) Cpu {
	if display == nil {
		display = NoDisplay{}
	}

	if keyboard == nil {
		keyboard = NoKeyboard{}
	}
	cpu := Cpu{
		Memory:   make([]uint8, memorySize),
		PC:       0x200,
		rng:      rngGenerator{},
		display:  display,
		keyboard: keyboard,
	}

	cpu.LoadInterpreter()

	return cpu
}
