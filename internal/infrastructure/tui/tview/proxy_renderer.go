package tview

import (
	"sync"
	"workercli/internal/domain/model"
	"workercli/internal/infrastructure/tui/tview/components"
	"workercli/pkg/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TViewProxyRenderer for ProxyCheck
type TViewProxyRenderer struct {
	boGhiNhatKy    *utils.Logger
	ketQua         *[]model.KetQuaTrungGian
	khoaKetQua     *sync.Mutex
	kenhKetQua     chan model.KetQuaTrungGian
	kenhDong       chan struct{}
	bangHienThi    *tview.Table
	ungDungHienThi *tview.Application
}

// NewTViewProxyRenderer creates a new TViewProxyRenderer
func NewTViewProxyRenderer(boGhiNhatKy *utils.Logger, ketQua *[]model.KetQuaTrungGian, khoaKetQua *sync.Mutex, kenhKetQua chan model.KetQuaTrungGian, kenhDong chan struct{}) *TViewProxyRenderer {
	return &TViewProxyRenderer{
		boGhiNhatKy: boGhiNhatKy,
		ketQua:      ketQua,
		khoaKetQua:  khoaKetQua,
		kenhKetQua:  kenhKetQua,
		kenhDong:    kenhDong,
	}
}

func (r *TViewProxyRenderer) Start() error {
	r.ungDungHienThi = tview.NewApplication()
	r.bangHienThi = tview.NewTable().SetBorders(true).SetFixed(1, 0).SetSelectable(true, false)

	r.ungDungHienThi.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			r.ungDungHienThi.Stop()
		}
		return event
	})

	go func() {
		hang := 1
		for {
			select {
			case ketQuaDon := <-r.kenhKetQua:
				r.khoaKetQua.Lock()
				*r.ketQua = append(*r.ketQua, ketQuaDon)
				components.RenderProxyTable(r.bangHienThi, r.ketQua, hang)
				hang++
				r.khoaKetQua.Unlock()
				r.ungDungHienThi.QueueUpdateDraw(func() {})
			case <-r.kenhDong:
				r.ungDungHienThi.QueueUpdateDraw(func() {})
				return
			}
		}
	}()

	if err := r.ungDungHienThi.SetRoot(r.bangHienThi, true).Run(); err != nil {
		r.boGhiNhatKy.Errorf("Không thể khởi động tview: %v", err)
		return err
	}
	return nil
}

func (r *TViewProxyRenderer) AddTaskResult(ketQua model.KetQua) {
	// Không được sử dụng trong renderer này
}

func (r *TViewProxyRenderer) AddProxyResult(ketQua model.KetQuaTrungGian) {
	select {
	case r.kenhKetQua <- ketQua:
		r.boGhiNhatKy.Infof("Đã thêm kết quả proxy vào kênh: %v", ketQua)
	case <-r.kenhDong:
		r.boGhiNhatKy.Info("Kênh proxy đã đóng, bỏ qua kết quả proxy")
	}
}

func (r *TViewProxyRenderer) Close() {
	// Đóng được xử lý trong usecase
}
