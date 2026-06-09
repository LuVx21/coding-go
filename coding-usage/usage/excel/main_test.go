package excel

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/xuri/excelize/v2"
)

var home = common_x.Home()

func Test_excel_00(t *testing.T) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	f.SetCellValue("Sheet1", cellIndex(1, "A"), "姓名")
	f.SetCellValue("Sheet1", cellIndex(1, "B"), "年龄")
	f.SetCellValue("Sheet1", cellIndex(2, "A"), "张三")
	f.SetCellValue("Sheet1", cellIndex(2, "B"), 28)
	f.SaveAs(home + "/demo.xlsx")

	cell, _ := f.GetCellValue("Sheet1", "A2")
	rows, _ := f.GetRows("Sheet1")
	fmt.Println(cell, rows)

}
func cellIndex(row int, col string) string { return col + strconv.Itoa(row) }

func Test_excel_01(t *testing.T) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	for idx, row := range [][]interface{}{
		{nil, "Apple", "Orange", "Pear"}, {"Small", 2, 3, 3},
		{"Normal", 5, 2, 4}, {"Large", 6, 7, 8},
	} {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetSheetRow("Sheet1", cell, &row)
	}
	if err := f.AddChart("Sheet1", "E1", &excelize.Chart{
		Type: excelize.Col3DClustered,
		Series: []excelize.ChartSeries{
			{
				Name:       "Sheet1!$A$2",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$2:$D$2",
			},
			{
				Name:       "Sheet1!$A$3",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$3:$D$3",
			},
			{
				Name:       "Sheet1!$A$4",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$4:$D$4",
			}},
		Title: []excelize.RichTextRun{
			{
				Text: "Fruit 3D Clustered Column Chart",
			},
		},
	},
	); err != nil {
		fmt.Println(err)
		return
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs(home + "/Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
