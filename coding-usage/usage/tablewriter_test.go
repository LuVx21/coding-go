package main

import (
	"os"
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

func Test_tablewriter_00(t *testing.T) {
	data := [][]any{
		{"Alice", 25, "New York"},
		{"Bob", 30, "Boston"},
	}

	symbols := tw.NewSymbolCustom("Nature").
		WithRow("~").
		WithColumn("|").
		WithTopLeft("🌱").
		WithTopMid("🌿").
		WithTopRight("🌱").
		WithMidLeft("🍃").
		WithCenter("❀").
		WithMidRight("🍃").
		WithBottomLeft("🌻").
		WithBottomMid("🌾").
		WithBottomRight("🌻")

	table := tablewriter.NewTable(os.Stdout, tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{Symbols: symbols})))
	table.Header("Name", "Age", "City")
	table.Bulk(data)
	table.Render()
}
