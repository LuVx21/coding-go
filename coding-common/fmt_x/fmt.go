package fmt_x

import (
    "github.com/jedib0t/go-pretty/v6/table"
    "os"
)

func Println(header []interface{}, rows ...[]interface{}) {
    tt := table.NewWriter()
    tt.SetOutputMirror(os.Stdout)
    tt.SetAutoIndex(true)
    tt.AppendHeader(header)
    for _, row := range rows {
        tt.AppendRow(row)
    }
    tt.Render()
}
