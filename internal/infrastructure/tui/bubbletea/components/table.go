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
func RenderTaskTable(ketQua *[]model.KetQua, selectedRow int) string {
	var buf strings.Builder

	header := []string{"Mã tác vụ", "Trạng thái"}
	headerWidths := []int{12, 20}
	headerRow := make([]string, len(header))
	for i, h := range header {
		headerRow[i] = lipgloss.NewStyle().Width(headerWidths[i]).Render(h)
	}
	buf.WriteString(headerStyle.Render(strings.Join(headerRow, "|")))
	buf.WriteString("\n")

	for i, kq := range *ketQua {
		row := []string{kq.MaTacVu, kq.TrangThai}
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
func RenderProxyTable(ketQua *[]model.KetQuaTrungGian, selectedRow int) string {
	var buf strings.Builder

	header := []string{"Giao thức", "Địa chỉ IP", "Cổng", "IP đã phân giải", "Trạng thái", "Lỗi"}
	headerWidths := []int{10, 15, 8, 15, 10, 20}
	headerRow := make([]string, len(header))
	for i, h := range header {
		headerRow[i] = lipgloss.NewStyle().Width(headerWidths[i]).Render(h)
	}
	buf.WriteString(headerStyle.Render(strings.Join(headerRow, "|")))
	buf.WriteString("\n")

	for i, kq := range *ketQua {
		row := []string{
			kq.TrungGian.GiaoDien,
			kq.TrungGian.DiaChi,
			kq.TrungGian.Cong,
			kq.DiaChi,
			kq.TrangThai,
			kq.LoiXayRa,
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
