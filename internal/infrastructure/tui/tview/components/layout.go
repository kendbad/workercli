package components

import (
	"workercli/internal/domain/model"

	"github.com/rivo/tview"
)

// RenderTaskTable renders the table for task results
func RenderTaskTable(tviewTable *tview.Table, results *[]model.Result, row int) {
	tviewTable.Clear()
	tviewTable.SetCell(0, 0, tview.NewTableCell("Task ID").SetAlign(tview.AlignLeft))
	tviewTable.SetCell(0, 1, tview.NewTableCell("Status").SetAlign(tview.AlignLeft))

	for i, result := range *results {
		tviewTable.SetCell(i+1, 0, tview.NewTableCell(result.TaskID).SetAlign(tview.AlignLeft))
		tviewTable.SetCell(i+1, 1, tview.NewTableCell(result.Status).SetAlign(tview.AlignLeft))
	}
}

// RenderProxyTable renders the table for proxy results
func RenderProxyTable(tviewTable *tview.Table, results *[]model.ProxyResult, row int) {
	tviewTable.Clear()
	tviewTable.SetCell(0, 0, tview.NewTableCell("Protocol").SetAlign(tview.AlignLeft))
	tviewTable.SetCell(0, 1, tview.NewTableCell("IP").SetAlign(tview.AlignLeft))
	tviewTable.SetCell(0, 2, tview.NewTableCell("Port").SetAlign(tview.AlignLeft))
	tviewTable.SetCell(0, 3, tview.NewTableCell("Resolved IP").SetAlign(tview.AlignLeft))
	tviewTable.SetCell(0, 4, tview.NewTableCell("Status").SetAlign(tview.AlignLeft))
	tviewTable.SetCell(0, 5, tview.NewTableCell("Error").SetAlign(tview.AlignLeft))

	for i, result := range *results {
		tviewTable.SetCell(i+1, 0, tview.NewTableCell(result.Proxy.Protocol).SetAlign(tview.AlignLeft))
		tviewTable.SetCell(i+1, 1, tview.NewTableCell(result.Proxy.IP).SetAlign(tview.AlignLeft))
		tviewTable.SetCell(i+1, 2, tview.NewTableCell(result.Proxy.Port).SetAlign(tview.AlignLeft))
		tviewTable.SetCell(i+1, 3, tview.NewTableCell(result.IP).SetAlign(tview.AlignLeft))
		tviewTable.SetCell(i+1, 4, tview.NewTableCell(result.Status).SetAlign(tview.AlignLeft))
		tviewTable.SetCell(i+1, 5, tview.NewTableCell(result.Error).SetAlign(tview.AlignLeft))
	}
}
