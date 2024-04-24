package fmt_x

import (
    "fmt"
    "github.com/jedib0t/go-pretty/v6/table"
    "os"
    "strconv"
    "strings"
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

//Println0
//+---+---+----+-----------------------------------------------------+
//| 1 | 1 | aa | 2024-05-08 11:34:52.202224 +0800 CST m=+0.000482917 |
//+---+---+----+-----------------------------------------------------+
func Println0[T any](rows ...[]T) {
    colNum := 0
    for _, _row := range rows {
        colNum = max(colNum, len(_row))
    }

    strRows := make([][]any, 0, len(rows))
    width := make([]int, colNum)
    for _, _row := range rows {
        strRow := make([]any, 0, colNum)
        for i := 0; i < colNum; i++ {
            var colStr string
            if i < len(_row) {
                col := _row[i]
                colStr = " " + fmt.Sprint(col) + " "
                width[i] = max(len(colStr), width[i])
            }
            strRow = append(strRow, colStr)
        }
        strRows = append(strRows, strRow)
    }

    format, line := "|", "+"
    for _, w := range width {
        format += "%" + strconv.Itoa(w) + "v|"
        line += strings.Repeat("-", w) + "+"
    }
    format += "\n"
    line += "\n"

    var sb strings.Builder
    sb.WriteString(line)
    for _, row := range strRows {
        sb.WriteString(fmt.Sprintf(format, row...))
        sb.WriteString(line)
    }
    fmt.Println(sb.String())
}
