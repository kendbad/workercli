package tui

import "workercli/internal/domain/model"

// Renderer defines the interface for TUI renderers
type Renderer interface {
	Start() error
	AddTaskResult(ketQua model.KetQua)
	AddProxyResult(ketQua model.KetQuaTrungGian)
	Close()
}
