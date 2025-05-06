package proxy

import (
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// Reader defines the interface for reading proxies.
type Reader interface {
	ReadProxies(nguon string) ([]model.Proxy, error)
}

// ProxyReader is the adapter for reading proxies.
type ProxyReader struct {
	boGhiNhatKy *utils.Logger
	boDoc       Reader // Specific implementation (e.g., file reader)
}

func NewProxyReader(boGhiNhatKy *utils.Logger, boDoc Reader) *ProxyReader {
	return &ProxyReader{
		boGhiNhatKy: boGhiNhatKy,
		boDoc:       boDoc,
	}
}

func (r *ProxyReader) ReadProxies(nguon string) ([]model.Proxy, error) {
	danhSachProxy, err := r.boDoc.ReadProxies(nguon)
	if err != nil {
		r.boGhiNhatKy.Errorf("Không thể đọc danh sách proxy từ %s: %v", nguon, err)
		return nil, err
	}
	r.boGhiNhatKy.Infof("Đã đọc %d proxy từ %s", len(danhSachProxy), nguon)
	return danhSachProxy, nil
}
