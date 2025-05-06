package proxy

import (
	"fmt"
	"workercli/internal/adapter/ipchecker"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

type BoKiemTraIP struct {
	boGhiNhatKy *utils.Logger
	boKiemTraIP ipchecker.IPChecker
}

func NewIPChecker(boGhiNhatKy *utils.Logger, loaiKetNoi string) *BoKiemTraIP {
	boKiemTraIP := ipchecker.NewIPChecker(loaiKetNoi, boGhiNhatKy)
	return &BoKiemTraIP{boGhiNhatKy: boGhiNhatKy, boKiemTraIP: boKiemTraIP}
}

func (c *BoKiemTraIP) KiemTraProxy(proxy model.Proxy, duongDanKiemTra string) (diaChi string, trangThai string, err error) {
	diaChi, maKetQua, err := c.boKiemTraIP.CheckIP(proxy, duongDanKiemTra)
	if err != nil {
		c.boGhiNhatKy.Errorf("Proxy %s://%s:%s thất bại: %v", proxy.GiaoDien, proxy.DiaChi, proxy.Cong, err)
		return "", fmt.Sprintf("Thất bại (%v)", err), err
	}

	if maKetQua != 200 {
		err = fmt.Errorf("mã trạng thái: %d", maKetQua)
		c.boGhiNhatKy.Errorf("Proxy %s://%s:%s trả về mã trạng thái: %d", proxy.GiaoDien, proxy.DiaChi, proxy.Cong, maKetQua)
		return "", fmt.Sprintf("Thất bại (mã trạng thái: %d)", maKetQua), err
	}

	return diaChi, "Thành công", nil
}
