package proxy

import (
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// BoKiemTra defines the interface for checking proxies.
type BoKiemTra interface {
	KiemTraTrungGian(trungGian model.TrungGian, duongDanKiemTra string) (diaChi string, trangThai string, err error)
}

// ProxyChecker is the adapter for checking proxies.
type ProxyChecker struct {
	boGhiNhatKy *utils.Logger
	boKiemTra   BoKiemTra // Specific implementation (e.g., IP checker)
}

func NewProxyChecker(boGhiNhatKy *utils.Logger, boKiemTra BoKiemTra) *ProxyChecker {
	return &ProxyChecker{
		boGhiNhatKy: boGhiNhatKy,
		boKiemTra:   boKiemTra,
	}
}

func (c *ProxyChecker) CheckProxy(trungGian model.TrungGian, duongDanKiemTra string) (diaChi string, trangThai string, err error) {
	diaChi, trangThai, err = c.boKiemTra.KiemTraTrungGian(trungGian, duongDanKiemTra)
	if err != nil {
		c.boGhiNhatKy.Errorf("Kiểm tra proxy thất bại %s://%s:%s: %v", trungGian.GiaoDien, trungGian.DiaChi, trungGian.Cong, err)
		return "", trangThai, err
	}
	c.boGhiNhatKy.Infof("Proxy %s://%s:%s trả về IP: %s", trungGian.GiaoDien, trungGian.DiaChi, trungGian.Cong, diaChi)
	return diaChi, trangThai, nil
}
