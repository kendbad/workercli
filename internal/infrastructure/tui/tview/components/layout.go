package components

import (
	"workercli/internal/domain/model"

	"github.com/rivo/tview"
)

// RenderTaskTable renders the table for task results
func RenderTaskTable(bangHienThi *tview.Table, ketQua *[]model.KetQua, hang int) {
	bangHienThi.Clear()
	bangHienThi.SetCell(0, 0, tview.NewTableCell("Mã tác vụ").SetAlign(tview.AlignLeft))
	bangHienThi.SetCell(0, 1, tview.NewTableCell("Trạng thái").SetAlign(tview.AlignLeft))

	for i, kq := range *ketQua {
		bangHienThi.SetCell(i+1, 0, tview.NewTableCell(kq.MaTacVu).SetAlign(tview.AlignLeft))
		bangHienThi.SetCell(i+1, 1, tview.NewTableCell(kq.TrangThai).SetAlign(tview.AlignLeft))
	}
}

// RenderProxyTable renders the table for proxy results
func RenderProxyTable(bangHienThi *tview.Table, ketQua *[]model.KetQuaTrungGian, hang int) {
	bangHienThi.Clear()
	bangHienThi.SetCell(0, 0, tview.NewTableCell("Giao thức").SetAlign(tview.AlignLeft))
	bangHienThi.SetCell(0, 1, tview.NewTableCell("Địa chỉ IP").SetAlign(tview.AlignLeft))
	bangHienThi.SetCell(0, 2, tview.NewTableCell("Cổng").SetAlign(tview.AlignLeft))
	bangHienThi.SetCell(0, 3, tview.NewTableCell("IP đã phân giải").SetAlign(tview.AlignLeft))
	bangHienThi.SetCell(0, 4, tview.NewTableCell("Trạng thái").SetAlign(tview.AlignLeft))
	bangHienThi.SetCell(0, 5, tview.NewTableCell("Lỗi").SetAlign(tview.AlignLeft))

	for i, kq := range *ketQua {
		bangHienThi.SetCell(i+1, 0, tview.NewTableCell(kq.TrungGian.GiaoDien).SetAlign(tview.AlignLeft))
		bangHienThi.SetCell(i+1, 1, tview.NewTableCell(kq.TrungGian.DiaChi).SetAlign(tview.AlignLeft))
		bangHienThi.SetCell(i+1, 2, tview.NewTableCell(kq.TrungGian.Cong).SetAlign(tview.AlignLeft))
		bangHienThi.SetCell(i+1, 3, tview.NewTableCell(kq.DiaChi).SetAlign(tview.AlignLeft))
		bangHienThi.SetCell(i+1, 4, tview.NewTableCell(kq.TrangThai).SetAlign(tview.AlignLeft))
		bangHienThi.SetCell(i+1, 5, tview.NewTableCell(kq.LoiXayRa).SetAlign(tview.AlignLeft))
	}
}
