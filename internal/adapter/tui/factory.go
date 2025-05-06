package tui

import (
	"sync"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// RendererFactory defines the interface for creating renderers
type RendererFactory interface {
	CreateTaskRenderer(boGhiNhatKy *utils.Logger, ketQua *[]model.KetQua, khoaKetQua *sync.Mutex, kenhKetQua chan model.KetQua, kenhDong chan struct{}) Renderer
	CreateProxyRenderer(boGhiNhatKy *utils.Logger, ketQua *[]model.KetQuaTrungGian, khoaKetQua *sync.Mutex, kenhKetQua chan model.KetQuaTrungGian, kenhDong chan struct{}) Renderer
}
