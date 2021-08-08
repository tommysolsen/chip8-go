package main

import (
	"chip8/src/chip8"
	"chip8/src/displays"
	"fmt"
	"os"
	"time"
)

func main() {
	display, err := displays.NewSDLRenderer(32)
	defer display.Dispose()
	if err != nil {
		display.Dispose()
		fmt.Println(err)
		os.Exit(-1)
	}
	defer display.Dispose()
	cpu := chip8.NewCPU(0xFFF, display, display)

	/*
		cpu.LoadProgram(bytes.NewReader([]byte{
			0x72, 0x05,
			0x73, 0x1B,
			0x61, 0x00,
			0xA0, 0x00,
			0xD0, 0x15, // Display 0
			0xD0, 0x35,

			// 0xF2, 0x1E, // Increment I by 5
			0x70, 0x05,

			0x12, 0x06,


		})) */

	f, err := os.Open("./games/BLINKY.ch8")

	err = cpu.LoadProgram(f)
	if err != nil {
		panic(err)
	}

	rate := int64(16)
	last := time.Now()

	// dumper := statedumpers.TableDumper{To: os.Stdout}

	memory := uint32(len(cpu.Memory))
	step := 0
	render := 0
	for true {
		cpu.Step()
		step++

		//dumper.DumpState(cpu)

		t := time.Since(last)
		if t.Milliseconds() > rate {
			fmt.Printf("Step: %015d\tRendered Frame:%010d\r", step, render)
			last = time.Now()
			cpu.DecrementTimers()
			render++
		}
		display.Render()
		//statedumpers.TableDumper{To: os.Stdout}.DumpState(cpu)
		if uint32(cpu.PC) > memory {
			break
		}
	}

	//dumper.DumpState(cpu)
}
