package fmt_x

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/luvx21/coding-go/coding-common/text_x"
)

const (
	// 格式: \033[显示方式;前景色;背景色m
	prefix      = "\033["
	colorRed    = prefix + "31m"
	colorGreen  = prefix + "32m"
	colorYellow = prefix + "33m"
	colorBlue   = prefix + "34m"
	colorPurple = prefix + "35m"
	colorCyan   = prefix + "36m"
	colorGray   = prefix + "37m"
	colorWhite  = prefix + "97m"
	colorReset  = prefix + "0m"
	// 红色文本绿色背景
	a = prefix + "31;42m"
	// 加粗红色
	b = prefix + "1;31m"
)

func PrintlnRow(item ...any) {
	Println(nil, item)
}

func Println(header []any, row ...[]any) {
	tt := table.NewWriter()
	tt.SetOutputMirror(os.Stdout)
	tt.SetAutoIndex(true)
	if header != nil {
		tt.AppendHeader(header)
	}
	for _, row := range row {
		tt.AppendRow(row)
	}
	tt.Render()
}

func PrintlnRow0(item ...any) {
	Println0(item)
}

// Println0
// +---+---+----+-----------------------------------------------------+
// | 1 | 1 | aa | 2024-05-08 11:34:52.202224 +0800 CST m=+0.000482917 |
// +---+---+----+-----------------------------------------------------+
func Println0[T any](rows ...[]T) {
	colNum := 0
	for _, _row := range rows {
		colNum = max(colNum, len(_row))
	}

	strRows := make([][]string, 0, len(rows))
	width := make([]int, colNum)
	for _, _row := range rows {
		strRow := make([]string, 0, colNum)
		for i := range colNum {
			var colStr string
			if i < len(_row) {
				col := _row[i]
				colStr = " " + fmt.Sprint(col) + " "
				width[i] = max(text_x.Width(colStr), width[i])
			}
			strRow = append(strRow, colStr)
		}
		strRows = append(strRows, strRow)
	}

	line := "+"
	for _, w := range width {
		line += strings.Repeat("-", w) + "+"
	}
	line += "\n"

	var sb strings.Builder
	sb.WriteString(line)
	for _, row := range strRows {
		sb.WriteString("|")
		for i, col := range row {
			for _ = range width[i] - text_x.Width(col) {
				sb.WriteString(" ")
			}
			sb.WriteString(col)
			sb.WriteString("|")
		}
		sb.WriteString("\n")
		sb.WriteString(line)
	}
	fmt.Println(sb.String())
}
