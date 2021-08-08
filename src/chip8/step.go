package chip8

import (
	"fmt"
	"strconv"
)

func (cpu *Cpu) NextInstruction() (uint8, uint8) {
	i1 := cpu.Memory[cpu.PC]
	i2 := cpu.Memory[cpu.PC+1]
	cpu.PC = cpu.PC + 2

	return i1, i2
}

func (cpu *Cpu) Step() {

	i1, i2 := cpu.NextInstruction()
	instruction := (uint16(i1) << 8) | uint16(i2)

	// JMP 1 instruction
	if instruction == 0x00E0 {
		cpu.display.Clear()
	} else if instruction == 0x00EE {
		cpu.SP -= 1
		cpu.PC = cpu.S[cpu.SP]
	} else if instruction >= 0x1000 && instruction < 0x2000 {
		cpu.PC = instruction - 0x1000
		return
	}

	// CALL instruction
	if instruction >= 0x2000 && instruction < 0x3000 {
		cpu.S[cpu.SP] = cpu.PC
		cpu.SP = cpu.SP + 1
		cpu.PC = instruction - 0x2000
		return
	}

	// SE Vx
	if instruction >= 0x3000 && instruction < 0x4000 {
		if cpu.V[i1-0x30] == i2 {
			cpu.PC = cpu.PC + 2
		}
		return
	}

	// SE Vx
	if instruction >= 0x4000 && instruction < 0x5000 {
		if cpu.V[i1-0x40] != i2 {
			cpu.PC = cpu.PC + 2
		}
		return
	}

	// SE Vx, Vy
	if instruction >= 0x5000 && instruction < 0x6000 {
		if cpu.V[i1-0x50] == cpu.V[i2>>4] {
			cpu.PC = cpu.PC + 2
		}
		return
	}

	// LD Vx instruction
	if instruction >= 0x6000 && instruction < 0x7000 {
		cpu.V[i1-0x60] = i2
		return
	}

	if instruction >= 0x7000 && instruction < 0x8000 {
		addr := i1 - 0x70
		cpu.V[addr] = uint8(cpu.V[addr] + i2)
		return
	}

	if instruction >= 0x8000 && instruction < 0x9000 && i2<<4 == 0 {
		cpu.V[i1-0x80] = cpu.V[i2>>4]
		return
	}

	if instruction >= 0x8000 && instruction < 0x9000 && (i2<<4 == 0x10) {
		cpu.V[i1-0x80] = cpu.V[i1-0x80] | cpu.V[i2>>4]
		return
	}

	if instruction >= 0x8000 && instruction < 0x9000 && (i2<<4 == 0x20) {
		cpu.V[i1-0x80] = cpu.V[i1-0x80] & cpu.V[i2>>4]
		return
	}
	if instruction >= 0x8000 && instruction < 0x9000 && (i2<<4 == 0x30) {
		cpu.V[i1-0x80] = cpu.V[i1-0x80] ^ cpu.V[i2>>4]
		return
	}
	if instruction >= 0x8000 && instruction < 0x9000 && (i2<<4 == 0x40) {
		result := uint16(cpu.V[i1-0x80]) + uint16(cpu.V[i2>>4])
		if result >= 0x100 {
			cpu.V[0x0F] = 1
		} else {
			cpu.V[0x0F] = 0
		}
		cpu.V[i1-0x80] = uint8(result)
		return
	}

	if instruction >= 0x8000 && instruction < 0x9000 && ((i2<<4 == 0x50) || (i2<<4 == 0x70)) {
		result := uint16(cpu.V[i1-0x80]) - uint16(cpu.V[i2>>4])

		var n uint8 = 0
		if i2<<4 == 0x70 {
			n = 1
		}

		if cpu.V[i1-0x80] > cpu.V[i2>>4] {

			cpu.V[0x0F] = 1 - n
		} else {
			cpu.V[0x0F] = 0 + n
		}
		cpu.V[i1-0x80] = uint8(result)
		return
	}

	if instruction >= 0x8000 && instruction < 0x9000 && (i2<<4 == 0x60) {
		x := cpu.V[i1-0x80]
		cpu.V[0x0F] = x & 0x01

		cpu.V[i1-0x80] = x >> 1
		return
	}

	if instruction >= 0x8000 && instruction < 0x9000 && (i2<<4 == 0xE0) {
		x := cpu.V[i1-0x80]
		if x&0x80 == 0x80 {
			cpu.V[0x0F] = 1
		} else {
			cpu.V[0x0F] = 0
		}

		cpu.V[i1-0x80] = x << 1
		return
	}

	if instruction >= 0x9000 && instruction < 0xA000 {
		if cpu.V[i1-0x90] != cpu.V[i2>>4] {
			cpu.PC = cpu.PC + 2
		}
		return
	}

	if instruction >= 0xA000 && instruction < 0xB000 {
		cpu.I = instruction - 0xA000
		return
	}

	if instruction >= 0xB000 && instruction < 0xC000 {
		cpu.PC = instruction - 0xB000 + uint16(cpu.V[0])
		return
	}

	if instruction >= 0xC000 && instruction < 0xD000 {
		cpu.V[i1-0xC0] = cpu.rng.GetRandom() & i2
		return
	}

	if instruction >= 0xD000 && instruction < 0xE000 {
		n := i2 & 0x0F
		cpu.display.SetSprite(cpu.V[i1-0xD0], cpu.V[i2>>4], cpu.Memory[cpu.I:cpu.I+uint16(n)])
		return
	}

	if instruction >= 0xF000 && i2 == 0x1E {
		cpu.I = cpu.I + uint16(cpu.V[i1-0xF0])
		return
	}

	if instruction >= 0xE09E && instruction <= 0xEF9E && i2 == 0x9E {
		if cpu.keyboard.IsDown(cpu.V[i1-0xE0]) {
			cpu.PC += 2
		}
		return
	}
	if instruction >= 0xE0A1 && instruction <= 0xEFA1 && i2 == 0xA1 {
		if cpu.keyboard.IsDown(cpu.V[i1-0xE0]) == false {
			cpu.PC += 2
		}
		return
	}

	if instruction >= 0xF000 && i2 == 0x07 {
		cpu.V[i1-0xF0] = cpu.DT
		return
	}

	if instruction > 0xF000 && i2 == 0x0A {
		cpu.V[i1-0xF0] = cpu.keyboard.WaitForKey()
		return
	}

	if instruction > 0xF000 && i2 == 0x15 {
		cpu.DT = cpu.V[i1-0xF0]
	}

	if instruction > 0xF000 && i2 == 0x18 {
		cpu.ST = cpu.V[i1-0xF0]
	}

	if instruction > 0xF000 && i2 == 0x29 {
		cpu.I = uint16(cpu.V[i1-0xF0]) * 4
	}

	if instruction > 0xF000 && i2 == 0x33 {
		s := fmt.Sprintf("%03d", cpu.V[i1-0xF0])

		for i, char := range s {
			i2, _ := strconv.ParseInt(strconv.QuoteRune(char)[1:2], 10, 10)
			cpu.Memory[cpu.I+uint16(i)] = uint8(i2)
		}
	}

	if instruction > 0xF00 && i2 == 0x55 {
		b := i1 - 0xF0
		for i := 0; uint8(i) <= b; i++ {
			cpu.Memory[cpu.I+uint16(i)] = cpu.V[i]
		}
	}

	if instruction > 0xF00 && i2 == 0x65 {
		b := i1 - 0xF0
		for i := 0; uint8(i) <= b; i++ {
			cpu.V[i] = cpu.Memory[cpu.I+uint16(i)]
		}
	}
}
