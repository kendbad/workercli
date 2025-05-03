package tui

import "workercli/internal/domain/model"

// Renderer defines the interface for TUI renderers
type Renderer interface {
	Start() error
	AddTaskResult(result model.Result)
	AddProxyResult(result model.ProxyResult)
	Close()
}
