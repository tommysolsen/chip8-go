package chip8

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func Test_Instruction_SYS(t *testing.T) {
	cpu := bootstrapTest([]byte{0x04, 0x00})
	before := cloneProcessor(cpu)
	cpu.Step()

	testProgramStepLength(cpu, before, 2, t)

}

func Test_Instruction_CLS(t *testing.T) {
	cpu := bootstrapTest([]byte{0x00, 0xE0})
	cpu.display = &TestClearDisplay{HasBeenCleared: false}

	cpu.Step()

	if cpu.display.(*TestClearDisplay).HasBeenCleared == false {
		t.Errorf("Screen has not been cleared")
	}
}

func Test_Instruction_RET(t *testing.T) {
	cpu := bootstrapTest([]byte{0x00, 0xEE})
	cpu.SP = 1
	cpu.S[0] = 0x123
	cpu.Step()
	if cpu.PC != 0x123 {
		t.Errorf("Expected PC to be 0x123, was: %#04x", cpu.PC)
	}
}

func Test_Instruction_JP_0(t *testing.T) {
	cpu := bootstrapTest([]byte{0x00, 0x01})
	before := cloneProcessor(cpu)
	testMemoryUnchanged(cpu, before, t)
}

func Test_Instruction_JP_1(t *testing.T) {
	cpu := bootstrapTest([]byte{0x11, 0x23})

	before := cloneProcessor(cpu)
	fmt.Println(cpu.PC, before.PC)
	cpu.Step()
	if cpu.PC != 0x123 {
		t.Errorf("Program Counter is not in the correct place. Expected 0x123, got %#04x", cpu.PC)
	}

	testMemoryUnchanged(cpu, before, t)

}

func Test_Instruction_CALL(t *testing.T) {
	cpu := bootstrapTest([]byte{0x22, 0x05, 0x00, 0x00, 0x00, 0xFF})
	cpu.Step()

	if cpu.PC != 0x205 {
		t.Errorf("Expected PC to be 0x205, was: %#04x", cpu.PC)
	}

	if cpu.SP != 1 {
		t.Errorf("Expected SP to be 1, was: %d", cpu.SP)
		t.FailNow()
	}

	if cpu.S[cpu.SP-1] != 0x202 {
		t.Errorf("Expected Top stack member to be 0x202 was %#03x", cpu.S[cpu.SP-1])
	}

}

func Test_Instruction_SE(t *testing.T) {
	cpuFails := bootstrapTest([]byte{
		0x30, 0xFF, 0x00, 0x00, 0x00, 0x00,
	})

	cpuPasses := cloneProcessor(cpuFails)

	cpuPasses.V[0] = 0xFF

	cpuPassesBefore := cloneProcessor(cpuPasses)
	cpuFailsBefore := cloneProcessor(cpuFails)

	cpuFails.Step()
	cpuPasses.Step()

	if cpuPasses.PC != 0x204 {
		t.Errorf("Expected passing cpu PC to be 0x204, was: %#03x", cpuPasses.PC)
	}

	if cpuFails.PC != 0x202 {
		t.Errorf("Expected failing cpu PC to be 0x202, was: %#03x", cpuFails.PC)
	}

	testMemoryUnchanged(cpuPasses, cpuPassesBefore, t)
	testMemoryUnchanged(cpuFails, cpuFailsBefore, t)
}

func Test_Instruction_SNE(t *testing.T) {
	cpuFails := bootstrapTest([]byte{
		0x40, 0xFF, 0x00, 0x00, 0x00, 0x00,
	})

	cpuPasses := cloneProcessor(cpuFails)
	cpuPasses.V[0] = 0xFF

	cpuPassesBefore := cloneProcessor(cpuPasses)
	cpuFailsBefore := cloneProcessor(cpuFails)

	cpuFails.Step()
	cpuPasses.Step()

	if cpuPasses.PC != 0x202 {
		t.Errorf("Expected passing cpu PC to be 0x204, was: %#03x", cpuPasses.PC)
	}

	if cpuFails.PC != 0x204 {
		t.Errorf("Expected failing cpu PC to be 0x202, was: %#03x", cpuFails.PC)
	}

	testMemoryUnchanged(cpuPasses, cpuPassesBefore, t)
	testMemoryUnchanged(cpuFails, cpuFailsBefore, t)
}

func Test_Instruction_SE_Vx_Vy(t *testing.T) {
	cpuPassing := bootstrapTest([]byte{
		0x50, 0x10, 0x00, 0x00,
	})
	cpuPassing.V[0] = 0x01
	cpuPassing.V[1] = 0x01

	cpuFailing := cloneProcessor(cpuPassing)
	cpuFailing.V[0] = 0x02
	cpuFailing.V[1] = 0x01

	cpuPassing.Step()
	cpuFailing.Step()

	if cpuPassing.PC != 0x204 {
		t.Errorf("Expected passing cpu PC to be 0x204, was: %#03x", cpuPassing.PC)
	}

	if cpuFailing.PC != 0x202 {
		t.Errorf("Expected failing cpu PC to be 0x202, was: %#03x", cpuFailing.PC)
	}
}

func Test_Instruction_K(t *testing.T) {
	cpu := bootstrapTest([]byte{0x60, 0x13, 0x61, 0xFF})
	before := cloneProcessor(cpu)
	cpu.Step()
	cpu.Step()

	testProgramStepLength(cpu, before, 4, t)

	if cpu.V[0] != 0x13 {
		t.Errorf("Expected 0x13 in V0, got %#02x", cpu.V[0])
	}

	if cpu.V[1] != 0xFF {
		t.Errorf("Expected 0xFF in V1, got %#02x", cpu.V[1])
	}

}

func Test_Instruction_ADD(t *testing.T) {
	cpu := bootstrapTest([]byte{0x70, 0x13, 0x61, 0xFF})
	cpu.V[0] = 0x10
	oldCpu := cloneProcessor(cpu)
	cpu.Step()

	if cpu.V[0] != 0x23 {
		t.Errorf("Expected V0 to be 0x23, was: %#02x", cpu.V[0])
	}
	testMemoryUnchanged(cpu, oldCpu, t)
}

func Test_Instruction_LD_Vx_Vy(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x10})
	before := cloneProcessor(cpu)
	cpu.V[0] = 0x00
	cpu.V[1] = 0xFF
	cpu.Step()

	if cpu.V[0] != cpu.V[1] {
		t.Errorf("Expected V0 to be equal to V1, 0 was %#02x and 1 was %#02x", cpu.V[0], cpu.V[1])
	}

	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_OR(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x11})
	cpu.V[0] = 0x02
	cpu.V[1] = 0x01

	before := cloneProcessor(cpu)
	cpu.Step()

	testArithmeticsInV0(cpu, before, "Bitwise OR", 0x3, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_AND(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x12})
	cpu.V[0] = 0xAA
	cpu.V[1] = 0x0F

	before := cloneProcessor(cpu)
	cpu.Step()

	testArithmeticsInV0(cpu, before, "Bitwise AND", 0x0A, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_XOR(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x13})
	cpu.V[0] = 0xAA
	cpu.V[1] = 0x0F

	before := cloneProcessor(cpu)
	cpu.Step()

	testArithmeticsInV0(cpu, before, "Bitwise XOR", 0xA5, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_ADD_VF_0(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x14})
	cpu.V[0] = 0xA0
	cpu.V[1] = 0x0F
	cpu.V[0x0F] = 0x01

	before := cloneProcessor(cpu)
	cpu.Step()

	testRestRegister(cpu, false, t)
	testArithmeticsInV0(cpu, before, "Add", 0xAF, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_ADD_VF_1(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x14})
	cpu.V[0] = 0xF0
	cpu.V[1] = 0x2F
	cpu.V[0x0F] = 0x10

	before := cloneProcessor(cpu)
	cpu.Step()

	testRestRegister(cpu, true, t)
	testArithmeticsInV0(cpu, before, "Add", 0x1F, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_SUB_VF_1(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x15})
	cpu.V[0] = 0xA0
	cpu.V[1] = 0x10
	cpu.V[0x0F] = 0x01

	before := cloneProcessor(cpu)
	cpu.Step()

	testRestRegister(cpu, true, t)
	testArithmeticsInV0(cpu, before, "Sub", 0x90, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_SUB_VF_0(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x15})
	cpu.V[0] = 0x10
	cpu.V[1] = 0x90
	cpu.V[0x0F] = 0x10

	before := cloneProcessor(cpu)
	cpu.Step()

	testRestRegister(cpu, false, t)
	testArithmeticsInV0(cpu, before, "Sub", 0x80, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_SHR_VF_1(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x16})
	cpu.V[0] = 0xFF
	cpu.V[0x0F] = 0x05

	before := cloneProcessor(cpu)
	cpu.Step()

	testRestRegister(cpu, true, t)
	testArithmeticsInV0(cpu, before, "Sub", 0x7F, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_SHR_VF_0(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x16})
	cpu.V[0] = 0x02
	cpu.V[0x0F] = 0x04

	before := cloneProcessor(cpu)
	cpu.Step()

	testRestRegister(cpu, false, t)
	testArithmeticsInV0(cpu, before, "SHR-1", 0x01, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_SUBN_VF_1(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x17})
	cpu.V[0] = 0xA0
	cpu.V[1] = 0x10
	cpu.V[0x0F] = 0x01

	before := cloneProcessor(cpu)
	cpu.Step()

	testRestRegister(cpu, false, t)
	testArithmeticsInV0(cpu, before, "Sub", 0x90, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_SUBN_VF_0(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x17})
	cpu.V[0] = 0x10
	cpu.V[1] = 0x90
	cpu.V[0x0F] = 0x10

	before := cloneProcessor(cpu)
	cpu.Step()

	testRestRegister(cpu, true, t)
	testArithmeticsInV0(cpu, before, "Sub", 0x80, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_SHL_VF_1(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x1E})
	cpu.V[0] = 0x81
	cpu.V[0x0F] = 0x05

	before := cloneProcessor(cpu)
	cpu.Step()

	testRestRegister(cpu, true, t)
	testArithmeticsInV0(cpu, before, "Sub", 0x02, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_Bitwise_SHL_VF_0(t *testing.T) {
	cpu := bootstrapTest([]byte{0x80, 0x1E})
	cpu.V[0] = 0x01
	cpu.V[0x0F] = 0x04

	before := cloneProcessor(cpu)
	cpu.Step()

	testRestRegister(cpu, false, t)
	testArithmeticsInV0(cpu, before, "SHR-1", 0x02, t)
	testMemoryUnchanged(cpu, before, t)
	testProgramStepLength(cpu, before, 2, t)
}

func Test_Instruction_SNE_Vx_Vy(t *testing.T) {
	cpu := bootstrapTest([]byte{0x90, 0x10, 0x90, 0x20})
	cpu.V[0] = 0x01
	cpu.V[1] = 0x01
	before := cloneProcessor(cpu)

	cpu.Step()
	testProgramStepLength(cpu, before, 2, t)

	cpu.Step()
	testProgramStepLength(cpu, before, 6, t)
}

func Test_Instruction_A(t *testing.T) {
	cpu := bootstrapTest([]byte{0xA1, 0x23})
	cpu.Step()

	if cpu.I != 0x0123 {
		t.Errorf("Expected I to be 0x123, was %#03x", cpu.I)
	}
}

func Test_Instruction_B(t *testing.T) {
	cpu := bootstrapTest([]byte{0xB1, 0x20})
	cpu.V[0] = 0x03
	cpu.Step()

	if cpu.PC != 0x0123 {
		t.Errorf("Expected I to be 0x123, was %#03x", cpu.PC)
	}

}

func Test_Instruction_C(t *testing.T) {
	cpu := bootstrapTest([]byte{0xC0, 0xF0})
	cpu.rng = rngGeneratorMock{Value: 0x55}
	before := cloneProcessor(cpu)
	cpu.Step()

	testArithmeticsInV0(cpu, before, "RNG using 01010101", 0x50, t)
	testProgramStepLength(cpu, before, 2, t)
	testMemoryUnchanged(cpu, before, t)

}

func Test_Instruction_ADD_I_Vx(t *testing.T) {
	cpu := bootstrapTest([]byte{0xF0, 0x1E})
	cpu.V[0] = 0x14

	cpu.Step()

	if cpu.I != uint16(cpu.V[0]) {
		t.Errorf("Expected I to be %#04x, was %#04x", cpu.V[0], cpu.I)
	}
}

func Test_Instruction_SKP_Vx(t *testing.T) {
	cpu := bootstrapTest([]byte{0xE0, 0x9E, 0xE1, 0x9E, 0x00, 0x00})
	cpu.V[1] = 1
	before := cloneProcessor(cpu)
	cpu.keyboard = &TestKeyboard{KeysDown: []uint8{1}}
	cpu.Step()
	testProgramStepLength(cpu, before, 2, t)
	cpu.Step()
	testProgramStepLength(cpu, before, 6, t)
}

func Test_Instruction_SKPNP_Vx(t *testing.T) {
	cpu := bootstrapTest([]byte{0xE0, 0xA1, 0xE1, 0xA1, 0x00, 0x00})
	cpu.V[1] = 1
	before := cloneProcessor(cpu)
	cpu.keyboard = &TestKeyboard{KeysDown: []uint8{1}}
	cpu.Step()
	testProgramStepLength(cpu, before, 4, t)
	cpu.Step()
	testProgramStepLength(cpu, before, 6, t)
}

func Test_Instruction_LD_Vx_DT(t *testing.T) {
	for i := 0; i < 16; i++ {
		cpu := bootstrapTest([]byte{0xF0 + uint8(i), 0x07})
		cpu.DT = 0x12
		cpu.Step()

		if cpu.V[i] != 0x12 {
			t.Errorf("Expected register %d to be 0x12: was %#02x", i, cpu.V[i])
		}
	}
}

func Test_Instruction_LD_Vx_K(t *testing.T) {
	for i := 0; i < 16; i++ {
		cpu := bootstrapTest([]byte{0xF0 + uint8(i), 0x0A})
		cpu.keyboard = &TestKeyboard{HasWaited: false, WaitFor: 0xFF}
		cpu.Step()

		if cpu.keyboard.(*TestKeyboard).HasWaited == false {
			t.Errorf("We have not waited for keypress")
		}

		if cpu.V[i] != 0xFF {
			t.Errorf("Expected V[%d] to be 0xFF: was: %#02x", i, cpu.V[i])
		}
	}
}

func Test_Instruction_LD_DT_Vx(t *testing.T) {
	for i := 0; i < 16; i++ {
		cpu := bootstrapTest([]byte{0xF0 + uint8(i), 0x15})
		r := uint8(rand.Uint32())
		cpu.V[i] = r
		cpu.Step()

		if cpu.DT != r {
			t.Errorf("Expected V[%d] to be %#02x: was: %#02x", i, r, cpu.DT)
		}
	}
}

func Test_Instruction_LD_ST_Vx(t *testing.T) {
	for i := 0; i < 16; i++ {
		cpu := bootstrapTest([]byte{0xF0 + uint8(i), 0x18})
		r := uint8(rand.Uint32())
		cpu.V[i] = r
		cpu.Step()

		if cpu.ST != r {
			t.Errorf("Expected V[%d] to be %#02x: was: %#02x", i, r, cpu.ST)
		}
	}
}

func Test_Instruction_LD_F_Vx(t *testing.T) {
	for i := 0; i < 16; i++ {
		cpu := bootstrapTest([]byte{0xF0 + uint8(i), 0x29})
		cpu.V[i] = uint8(i)
		cpu.Step()

		if cpu.I != uint16(i*4) {
			t.Errorf("Expected I to be %#04x: was: %#04x", i*4, cpu.I)
		}
	}
}

func Test_Instruction_LD_B_Vx(t *testing.T) {
	for i := 0; i < 16; i++ {
		cpu := bootstrapTest([]byte{0xF0 + uint8(i), 0x33})
		cpu.V[i] = 123
		cpu.I = 0x123

		cpu.Step()

		if cpu.Memory[0x123] != 0x01 {
			t.Errorf("Using register %d, Expected hundreds field to be 1, got %d", i, cpu.Memory[0x123])
		}

		if cpu.Memory[0x123] != 0x01 {
			t.Errorf("Using register %d, Expected tens field to be 2 got %d", i, cpu.Memory[0x124])
		}

		if cpu.Memory[0x123] != 0x01 {
			t.Errorf("Using register %d, Expected decimal field to be 3, got %d", i, cpu.Memory[0x125])
		}

	}
}

func Test_Instruction_LD_I_Vx(t *testing.T) {
	for i := 0; i < 16; i++ {
		cpu := bootstrapTest([]byte{0xF0 + uint8(i), 0x55})
		before := cloneProcessor(cpu)
		for i2 := 0; i2 <= i; i2++ {
			cpu.V[i2] = uint8(i2)
		}
		cpu.I = 0x123

		cpu.Step()

		for i2 := 0; i2 <= i; i2++ {
			if cpu.Memory[cpu.I+uint16(i2)] != uint8(i2) {
				t.Errorf("In run %d, Expected memory %#04x to be %#02x: was %#02x ", i, cpu.I+uint16(i2), i2, cpu.Memory[cpu.I+uint16(i2)])
			}
		}

		if cpu.Memory[cpu.I+uint16(i)+1] != before.Memory[cpu.I+uint16(i)+1] {
			t.Errorf("The memory address %#04x should be untouched", cpu.I+uint16(i)+1)
		}

	}
}

func Test_Instruction_LD_VX_I(t *testing.T) {
	for i := 0; i < 16; i++ {
		cpu := bootstrapTest([]byte{0xF0 + uint8(i), 0x65})
		for i2 := 0; i2 <= i; i2++ {
			cpu.Memory[0x123+i2] = uint8(i2)
		}
		cpu.I = 0x123

		cpu.Step()

		for i2 := 0; i2 <= i; i2++ {
			if cpu.V[uint16(i2)] != uint8(i2) {
				t.Errorf("In run %d, Expected V[%d] to be %#02x: was %#02x ", i, i2, i2, cpu.Memory[cpu.I+uint16(i2)])
			}
		}
	}
}

func bootstrapTest(code []byte) Cpu {
	cpu := Cpu{
		Memory:  make([]uint8, 0x200+len(code)),
		PC:      0x200,
		rng:     rngGenerator{},
		display: NoDisplay{},
	}
	err := cpu.LoadProgram(bytes.NewReader(code))
	if err != nil {
		panic(err.Error())
	}
	return cpu
}

func testRestRegister(cpu Cpu, value bool, t *testing.T) {
	var v uint8 = 0
	if value {
		v = 1
	}
	if cpu.V[0x0F] != v {
		t.Errorf("Expected Register F to be %#02x, was %#02x", v, cpu.V[0x0F])
	}

}
func testArithmeticsInV0(cpu Cpu, before Cpu, label string, value uint8, t *testing.T) {
	if cpu.V[0] != value {
		t.Errorf("%s failed: expected %08b from %08b and %08b. Got: %08b", label, value, before.V[0], before.V[1], cpu.V[0])
	}
}

func testProgramStepLength(cpu Cpu, oldCpu Cpu, l int, t *testing.T) {
	if cpu.PC != oldCpu.PC+uint16(l) {
		t.Errorf("Expected a step count of %d, got %d", l, int64(math.Abs(float64(cpu.PC)-float64(oldCpu.PC))))
	}
}

func testMemoryUnchanged(cpu Cpu, oldCpu Cpu, t *testing.T) {
	for i, val := range cpu.Memory {
		if val != oldCpu.Memory[i] {
			t.Errorf("Memory inconsistent expected %d in %#04x, got %d", oldCpu.Memory[i], i, val)
		}
	}
}

func cloneProcessor(cpu Cpu) Cpu {
	nCpu := Cpu{
		Memory: make([]uint8, len(cpu.Memory)),
	}

	for i, v := range cpu.Memory {
		nCpu.Memory[i] = v
	}

	for i, v := range cpu.S {
		nCpu.S[i] = v
	}

	for i, v := range cpu.V {
		nCpu.V[i] = v
	}

	nCpu.PC = cpu.PC
	nCpu.DT = cpu.DT
	nCpu.ST = cpu.ST
	nCpu.SP = cpu.SP

	return nCpu
}
