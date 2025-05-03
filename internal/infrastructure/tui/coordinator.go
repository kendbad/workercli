package tui

import (
	"sync"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// TUIRenderer định nghĩa giao diện trừu tượng cho renderer
type TUIRenderer interface {
	Start() error
	AddTaskResult(result model.Result)
	AddProxyResult(result model.ProxyResult)
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
	logger   *utils.Logger
	mode     TUIMode
	renderer TUIRenderer
	wg       sync.WaitGroup
}

// NewTUIUseCase khởi tạo TUIUseCase
func NewTUIUseCase(logger *utils.Logger, mode string, renderer TUIRenderer) *TUIUseCase {
	return &TUIUseCase{
		logger:   logger,
		mode:     validateMode(mode),
		renderer: renderer,
	}
}

// Start khởi động TUI
func (uc *TUIUseCase) Start() error {
	if uc.renderer == nil {
		return nil
	}
	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		uc.logger.Info("Starting TUI")
		if err := uc.renderer.Start(); err != nil {
			uc.logger.Errorf("Failed to start TUI: %v", err)
		}
	}()
	return nil
}

// AddTaskResult thêm kết quả task
func (uc *TUIUseCase) AddTaskResult(result model.Result) {
	if uc.renderer != nil {
		uc.renderer.AddTaskResult(result)
	}
}

// AddProxyResult thêm kết quả proxy
func (uc *TUIUseCase) AddProxyResult(result model.ProxyResult) {
	if uc.renderer != nil {
		uc.renderer.AddProxyResult(result)
	}
}

// Close đóng TUI
func (uc *TUIUseCase) Close() {
	uc.wg.Wait()

	if uc.renderer != nil {
		uc.logger.Info("Closing TUI")
		uc.renderer.Close()
	}
}

// validateTUIMode xác thực TUIMode
func validateMode(mode string) TUIMode {
	switch mode {
	case string(TUIModeTView), string(TUIModeBubbleTea):
		return TUIMode(mode)
	default:
		return TUIModeTView
	}
}
