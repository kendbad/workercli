package proxy

import (
	"fmt"
	"strings"
	"workercli/internal/domain/model"
)

// Proxy represents a proxy configuration
type Proxy struct {
	GiaoDien string
	DiaChi   string
	Cong     string
}

// ParseProxy parses a proxy string (e.g., "http://1.2.3.4:8080") into a Proxy struct
func ParseProxy(proxyStr string) (model.Proxy, error) {
	parts := strings.SplitN(proxyStr, "://", 2)
	if len(parts) != 2 {
		return model.Proxy{}, fmt.Errorf("định dạng proxy không hợp lệ: %s", proxyStr)
	}

	addrParts := strings.Split(parts[1], ":")
	if len(addrParts) != 2 {
		return model.Proxy{}, fmt.Errorf("địa chỉ proxy không hợp lệ: %s", parts[1])
	}
	return model.Proxy{GiaoDien: parts[0], DiaChi: addrParts[0], Cong: addrParts[1]}, nil
}
