package tview

import (
	"sync"
	"workercli/internal/domain/model"
	"workercli/internal/infrastructure/tui/tview/components"
	"workercli/pkg/utils"

	"github.com/rivo/tview"
)

// TViewRenderer for BatchTask
type TViewRenderer struct {
	boGhiNhatKy    *utils.Logger
	ketQua         *[]model.KetQua
	khoaKetQua     *sync.Mutex
	kenhKetQua     chan model.KetQua
	kenhDong       chan struct{}
	bangHienThi    *tview.Table
	ungDungHienThi *tview.Application
}

// NewTViewRenderer creates a new TViewRenderer
func NewTViewRenderer(boGhiNhatKy *utils.Logger, ketQua *[]model.KetQua, khoaKetQua *sync.Mutex, kenhKetQua chan model.KetQua, kenhDong chan struct{}) *TViewRenderer {
	return &TViewRenderer{
		boGhiNhatKy: boGhiNhatKy,
		ketQua:      ketQua,
		khoaKetQua:  khoaKetQua,
		kenhKetQua:  kenhKetQua,
		kenhDong:    kenhDong,
	}
}

func (r *TViewRenderer) Start() error {
	r.ungDungHienThi = tview.NewApplication()
	r.bangHienThi = tview.NewTable().SetBorders(true).SetFixed(1, 0).SetSelectable(true, false)

	go func() {
		hang := 1
		for {
			select {
			case ketQuaDon := <-r.kenhKetQua:
				r.khoaKetQua.Lock()
				*r.ketQua = append(*r.ketQua, ketQuaDon)
				components.RenderTaskTable(r.bangHienThi, r.ketQua, hang)
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

func (r *TViewRenderer) AddTaskResult(ketQua model.KetQua) {
	select {
	case r.kenhKetQua <- ketQua:
		r.boGhiNhatKy.Infof("Đã thêm kết quả tác vụ vào kênh: %v", ketQua)
	case <-r.kenhDong:
		r.boGhiNhatKy.Info("Kênh tác vụ đã đóng, bỏ qua kết quả")
	}
}

func (r *TViewRenderer) AddProxyResult(ketQua model.KetQuaProxy) {
	// Không được sử dụng trong renderer này
}

func (r *TViewRenderer) Close() {
	// Đóng được xử lý trong usecase
}
