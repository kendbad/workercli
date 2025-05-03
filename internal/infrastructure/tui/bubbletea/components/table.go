package components

import (
	"strings"
	"workercli/internal/domain/model"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle      = lipgloss.NewStyle().Bold(true)
	normalRowStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	selectedRowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("105"))
)

// RenderTaskTable renders the table for task results
func RenderTaskTable(results *[]model.Result, selectedRow int) string {
	var buf strings.Builder

	header := []string{"Task ID", "Status"}
	headerWidths := []int{12, 20}
	headerRow := make([]string, len(header))
	for i, h := range header {
		headerRow[i] = lipgloss.NewStyle().Width(headerWidths[i]).Render(h)
	}
	buf.WriteString(headerStyle.Render(strings.Join(headerRow, "|")))
	buf.WriteString("\n")

	for i, res := range *results {
		row := []string{res.TaskID, res.Status}
		formattedRow := make([]string, len(row))
		for j, cell := range row {
			formattedRow[j] = lipgloss.NewStyle().Width(headerWidths[j]).Render(cell)
		}
		line := strings.Join(formattedRow, "|")
		if i == selectedRow {
			buf.WriteString(selectedRowStyle.Render(line))
		} else {
			buf.WriteString(normalRowStyle.Render(line))
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

// RenderProxyTable renders the table for proxy results
func RenderProxyTable(results *[]model.ProxyResult, selectedRow int) string {
	var buf strings.Builder

	header := []string{"Protocol", "IP", "Port", "Resolved IP", "Status", "Error"}
	headerWidths := []int{10, 15, 8, 15, 10, 20}
	headerRow := make([]string, len(header))
	for i, h := range header {
		headerRow[i] = lipgloss.NewStyle().Width(headerWidths[i]).Render(h)
	}
	buf.WriteString(headerStyle.Render(strings.Join(headerRow, "|")))
	buf.WriteString("\n")

	for i, res := range *results {
		row := []string{
			res.Proxy.Protocol,
			res.Proxy.IP,
			res.Proxy.Port,
			res.IP,
			res.Status,
			res.Error,
		}
		formattedRow := make([]string, len(row))
		for j, cell := range row {
			formattedRow[j] = lipgloss.NewStyle().Width(headerWidths[j]).Render(cell)
		}
		line := strings.Join(formattedRow, "|")
		if i == selectedRow {
			buf.WriteString(selectedRowStyle.Render(line))
		} else {
			buf.WriteString(normalRowStyle.Render(line))
		}
		buf.WriteString("\n")
	}

	return buf.String()
}
