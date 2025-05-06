package proxy

import (
	"fmt"
	"workercli/internal/domain/model"

	"workercli/pkg/utils"
)

type Checker struct {
	boGhiNhatKy *utils.Logger
	boKiemTraIP *BoKiemTraIP
}

func NewChecker(boGhiNhatKy *utils.Logger, loaiKetNoi string) *Checker {
	boKiemTraIP := NewIPChecker(boGhiNhatKy, loaiKetNoi)
	return &Checker{boGhiNhatKy: boGhiNhatKy, boKiemTraIP: boKiemTraIP}
}

func (c *Checker) CheckProxy(trungGian model.TrungGian, duongDanKiemTra string) (diaChi string, trangThai string, err error) {
	diaChi, trangThaiKetQua, err := c.boKiemTraIP.KiemTraTrungGian(trungGian, duongDanKiemTra)
	if err != nil {
		c.boGhiNhatKy.Errorf("Proxy %s://%s:%s thất bại: %v", trungGian.GiaoDien, trungGian.DiaChi, trungGian.Cong, err)
		return "", fmt.Sprintf("Thất bại (%v)", err), err
	}

	// Kiểm tra nếu trangThaiKetQua là "Thành công" thì trả về kết quả thành công
	if trangThaiKetQua != "Thành công" {
		c.boGhiNhatKy.Errorf("Proxy %s://%s:%s trả về trạng thái: %s", trungGian.GiaoDien, trungGian.DiaChi, trungGian.Cong, trangThaiKetQua)
		return "", trangThaiKetQua, fmt.Errorf("kiểm tra proxy không thành công: %s", trangThaiKetQua)
	}

	c.boGhiNhatKy.Infof("Proxy %s://%s:%s trả về IP: %s", trungGian.GiaoDien, trungGian.DiaChi, trungGian.Cong, diaChi)
	return diaChi, "Thành công", nil
}
