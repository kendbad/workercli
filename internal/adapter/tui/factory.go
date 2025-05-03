package tui

import (
	"sync"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// RendererFactory defines the interface for creating renderers
type RendererFactory interface {
	CreateTaskRenderer(logger *utils.Logger, results *[]model.Result, resultsMu *sync.Mutex, resultChan chan model.Result, closeChan chan struct{}) Renderer
	CreateProxyRenderer(logger *utils.Logger, results *[]model.ProxyResult, resultsMu *sync.Mutex, resultChan chan model.ProxyResult, closeChan chan struct{}) Renderer
}
