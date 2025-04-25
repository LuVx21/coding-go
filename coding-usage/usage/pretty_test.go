package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func Test_p_01(tt *testing.T) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	header := table.Row{"ID", "IP", "Num", "PacketsRecv", "PacketLoss", "AvgRtt"}
	t.AppendHeader(header)
	t.SetTitle("汇总")
	t.SetAutoIndex(true)
	for i := 1; i <= 5; i++ {
		row := table.Row{i, fmt.Sprintf("10.0.0.%v", i), i + 4, i, i, "AppendRow"}
		t.AppendRow(row)
	}
	t.AppendSeparator()
	t.AppendFooter(table.Row{"", "", "", "", "Total", 10000})
	//t.SetStyle(table.StyleLight)
	//t.SetStyle(table.StyleColoredBright)
	t.Render()
}

func Test_merge_00(tt *testing.T) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.SetTitle("所有模型")
	t.AppendHeader(table.Row{"No", "模型", "服务商", "No", "模型", "服务商"})
	t.AppendRow(table.Row{}, rowConfigAutoMerge)
	t.AppendRows([]table.Row{
		{1, 2, 3, 4, 5, 6},
		{11, 22, 3, 44, 55, 6},
	}, rowConfigAutoMerge)

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 3, AutoMerge: true, Align: text.AlignRight},
		{Number: 6, AutoMerge: true, Align: text.AlignRight},
	})
	t.Style().Options.SeparateRows = true
	t.Render()
}

func Test_custom_00(tt *testing.T) {
	leftUp, upMid, rightUp := "╭", "┬", "╮"
	leftMid, midMid, rightMid := "├", "┼", "┤"
	leftDown, downMid, rightDown := "╰", "┴", "╯"
	split, line := "│", "─"

	colWidth := 4
	s, ctx := strings.Repeat(line, colWidth), strings.Repeat(" ", colWidth)
	fmt.Println(leftUp + s + upMid + s + rightUp)
	fmt.Println(split + ctx + split + ctx + split)
	fmt.Println(leftMid + s + midMid + s + rightMid)
	fmt.Println(split + ctx + split + ctx + split)
	fmt.Println(split + ctx + split + ctx + split)
	fmt.Println(leftDown + s + downMid + s + rightDown)
}
