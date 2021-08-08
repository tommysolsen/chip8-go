package statedumpers

import (
	"chip8/src/chip8"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"strconv"
)

type TableDumper struct {
	To io.Writer
}

func (t TableDumper) DumpState(c chip8.Cpu) {
	fmt.Println("")
	table := tablewriter.NewWriter(t.To)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"Subject", "Value(Hex)", "Value(int)", "RValue(Hex)", "RValue(int)"})

	var data [][]string

	for i, v := range c.V {
		data = append(
			data,
			[]string{
				fmt.Sprintf("V%d", i),
				fmt.Sprintf("%#02x", v),
				strconv.FormatInt(int64(v), 10),
				"---",
				"---",
			})
	}

	data = append(
		data,
		[]string{
			"I",
			fmt.Sprintf("%#02x", c.I),
			strconv.FormatInt(int64(c.I), 10),
			fmt.Sprintf("%#02x", c.Memory[c.I]),
			strconv.FormatInt(int64(c.Memory[c.I]), 10),
		})

	data = append(
		data,
		[]string{
			"DT",
			"0x" + strconv.FormatInt(int64(c.DT), 16),
			strconv.FormatInt(int64(c.DT), 10),
			"---",
			"---",
		})
	data = append(
		data,
		[]string{
			"ST",
			fmt.Sprintf("%#02x", c.ST),
			strconv.FormatInt(int64(c.ST), 10),
			"---",
			"---",
		})

	data = append(
		data,
		[]string{
			"PC",
			fmt.Sprintf("%#02x", c.PC),
			strconv.FormatInt(int64(c.PC), 10),
			fmt.Sprintf("%#02x", c.Memory[c.PC]),
			strconv.FormatInt(int64(c.Memory[c.PC]), 10),
		})

	data = append(
		data,
		[]string{
			"PC+1",
			fmt.Sprintf("%#02x", c.PC),
			strconv.FormatInt(int64(c.PC+1), 10),
			fmt.Sprintf("%#02x", c.Memory[c.PC+1]),
			strconv.FormatInt(int64(c.Memory[c.PC+1]), 10),
		})

	data = append(
		data,
		[]string{
			"SP",
			"0x" + strconv.FormatInt(int64(c.SP), 16),
			strconv.FormatInt(int64(c.SP), 10),
			"0x" + strconv.FormatInt(int64(c.S[c.SP]), 16),
			strconv.FormatInt(int64(c.S[c.SP]), 10),
		})

	table.AppendBulk(data)
	table.Render()
}
