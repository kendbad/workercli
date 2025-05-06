package tui

import (
	"os"
	"sync"
	"workercli/internal/adapter/tui"
	"workercli/internal/domain/model"
	"workercli/internal/infrastructure/tui/bubbletea"
	"workercli/internal/infrastructure/tui/tview"
	"workercli/pkg/utils"

	"golang.org/x/term"
)

// RendererFactoryImpl implements RendererFactory
type RendererFactoryImpl struct {
	boGhiNhatKy *utils.Logger
	kieuHienThi tui.TUIMode
}

// NewRendererFactory creates a new factory
func NewRendererFactory(boGhiNhatKy *utils.Logger, kieuHienThi string) tui.RendererFactory {
	kieuHienThiHopLe := validateTUIMode(kieuHienThi)
	return &RendererFactoryImpl{
		boGhiNhatKy: boGhiNhatKy,
		kieuHienThi: kieuHienThiHopLe,
	}
}

func (f *RendererFactoryImpl) CreateTaskRenderer(boGhiNhatKy *utils.Logger, ketQua *[]model.KetQua, khoaKetQua *sync.Mutex, kenhKetQua chan model.KetQua, kenhDong chan struct{}) tui.Renderer {
	switch f.kieuHienThi {
	case tui.TUIModeTView:
		return tview.NewTViewRenderer(boGhiNhatKy, ketQua, khoaKetQua, kenhKetQua, kenhDong)
	case tui.TUIModeBubbleTea:
		return bubbletea.NewBubbleTeaRenderer(boGhiNhatKy, ketQua, khoaKetQua, kenhKetQua, kenhDong)
	default:
		return nil
	}
}

func (f *RendererFactoryImpl) CreateProxyRenderer(boGhiNhatKy *utils.Logger, ketQua *[]model.KetQuaTrungGian, khoaKetQua *sync.Mutex, kenhKetQua chan model.KetQuaTrungGian, kenhDong chan struct{}) tui.Renderer {
	switch f.kieuHienThi {
	case tui.TUIModeTView:
		return tview.NewTViewProxyRenderer(boGhiNhatKy, ketQua, khoaKetQua, kenhKetQua, kenhDong)
	case tui.TUIModeBubbleTea:
		return bubbletea.NewBubbleTeaProxyRenderer(boGhiNhatKy, ketQua, khoaKetQua, kenhKetQua, kenhDong)
	default:
		return nil
	}
}

// validateTUIMode validates TUIMode
func validateTUIMode(kieuHienThi string) tui.TUIMode {
	switch kieuHienThi {
	case string(tui.TUIModeTView), string(tui.TUIModeBubbleTea):
		if !hasTTY() {
			return tui.TUIModeNone
		}
		return tui.TUIMode(kieuHienThi)
	default:
		return tui.TUIModeTView
	}
}

// hasTTY checks for TTY
func hasTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
