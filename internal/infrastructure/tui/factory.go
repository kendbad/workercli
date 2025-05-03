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
	logger *utils.Logger
	mode   tui.TUIMode
}

// NewRendererFactory creates a new factory
func NewRendererFactory(logger *utils.Logger, mode string) tui.RendererFactory {
	validatedMode := validateTUIMode(mode)
	return &RendererFactoryImpl{
		logger: logger,
		mode:   validatedMode,
	}
}

func (f *RendererFactoryImpl) CreateTaskRenderer(logger *utils.Logger, results *[]model.Result, resultsMu *sync.Mutex, resultChan chan model.Result, closeChan chan struct{}) tui.Renderer {
	switch f.mode {
	case tui.TUIModeTView:
		return tview.NewTViewRenderer(logger, results, resultsMu, resultChan, closeChan)
	case tui.TUIModeBubbleTea:
		return bubbletea.NewBubbleTeaRenderer(logger, results, resultsMu, resultChan, closeChan)
	default:
		return nil
	}
}

func (f *RendererFactoryImpl) CreateProxyRenderer(logger *utils.Logger, results *[]model.ProxyResult, resultsMu *sync.Mutex, resultChan chan model.ProxyResult, closeChan chan struct{}) tui.Renderer {
	switch f.mode {
	case tui.TUIModeTView:
		return tview.NewTViewProxyRenderer(logger, results, resultsMu, resultChan, closeChan)
	case tui.TUIModeBubbleTea:
		return bubbletea.NewBubbleTeaProxyRenderer(logger, results, resultsMu, resultChan, closeChan)
	default:
		return nil
	}
}

// validateTUIMode validates TUIMode
func validateTUIMode(mode string) tui.TUIMode {
	switch mode {
	case string(tui.TUIModeTView), string(tui.TUIModeBubbleTea):
		if !hasTTY() {
			return tui.TUIModeNone
		}
		return tui.TUIMode(mode)
	default:
		return tui.TUIModeTView
	}
}

// hasTTY checks for TTY
func hasTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
