package main

import (
    "fmt"
    "github.com/jedib0t/go-pretty/v6/table"
    "os"
    "testing"
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
