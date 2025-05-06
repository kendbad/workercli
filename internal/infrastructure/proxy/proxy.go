package proxy

import (
	"fmt"
	"strings"
	"workercli/internal/domain/model"
)

// TrungGian represents a proxy configuration
type TrungGian struct {
	GiaoDien string
	DiaChi   string
	Cong     string
}

// ParseProxy parses a proxy string (e.g., "http://1.2.3.4:8080") into a TrungGian struct
func ParseProxy(proxyStr string) (model.TrungGian, error) {
	parts := strings.SplitN(proxyStr, "://", 2)
	if len(parts) != 2 {
		return model.TrungGian{}, fmt.Errorf("định dạng proxy không hợp lệ: %s", proxyStr)
	}

	addrParts := strings.Split(parts[1], ":")
	if len(addrParts) != 2 {
		return model.TrungGian{}, fmt.Errorf("địa chỉ proxy không hợp lệ: %s", parts[1])
	}
	return model.TrungGian{GiaoDien: parts[0], DiaChi: addrParts[0], Cong: addrParts[1]}, nil
}
