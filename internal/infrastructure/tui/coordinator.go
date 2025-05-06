package tui

import (
	"sync"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// TUIRenderer định nghĩa giao diện trừu tượng cho renderer
type TUIRenderer interface {
	Start() error
	AddTaskResult(ketQua model.KetQua)
	AddProxyResult(ketQua model.KetQuaTrungGian)
	Close()
}

// TUIMode định nghĩa loại giao diện
type TUIMode string

const (
	TUIModeNone      TUIMode = ""
	TUIModeTView     TUIMode = "tview"
	TUIModeBubbleTea TUIMode = "bubbletea"
)

// TUIUseCase quản lý logic của TUI
type TUIUseCase struct {
	boGhiNhatKy *utils.Logger
	kieuHienThi TUIMode
	boHienThi   TUIRenderer
	nhomCho     sync.WaitGroup
}

// NewTUIUseCase khởi tạo TUIUseCase
func NewTUIUseCase(boGhiNhatKy *utils.Logger, kieuHienThi string, boHienThi TUIRenderer) *TUIUseCase {
	return &TUIUseCase{
		boGhiNhatKy: boGhiNhatKy,
		kieuHienThi: validateMode(kieuHienThi),
		boHienThi:   boHienThi,
	}
}

// Start khởi động TUI
func (uc *TUIUseCase) Start() error {
	if uc.boHienThi == nil {
		return nil
	}
	uc.nhomCho.Add(1)
	go func() {
		defer uc.nhomCho.Done()
		uc.boGhiNhatKy.Info("Đang khởi động TUI")
		if err := uc.boHienThi.Start(); err != nil {
			uc.boGhiNhatKy.Errorf("Không thể khởi động TUI: %v", err)
		}
	}()
	return nil
}

// AddTaskResult thêm kết quả task
func (uc *TUIUseCase) AddTaskResult(ketQua model.KetQua) {
	if uc.boHienThi != nil {
		uc.boHienThi.AddTaskResult(ketQua)
	}
}

// AddProxyResult thêm kết quả proxy
func (uc *TUIUseCase) AddProxyResult(ketQua model.KetQuaTrungGian) {
	if uc.boHienThi != nil {
		uc.boHienThi.AddProxyResult(ketQua)
	}
}

// Close đóng TUI
func (uc *TUIUseCase) Close() {
	uc.nhomCho.Wait()

	if uc.boHienThi != nil {
		uc.boGhiNhatKy.Info("Đang đóng TUI")
		uc.boHienThi.Close()
	}
}

// validateTUIMode xác thực TUIMode
func validateMode(kieuHienThi string) TUIMode {
	switch kieuHienThi {
	case string(TUIModeTView), string(TUIModeBubbleTea):
		return TUIMode(kieuHienThi)
	default:
		return TUIModeTView
	}
}
