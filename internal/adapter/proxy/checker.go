package proxy

import (
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// BoKiemTra định nghĩa giao diện để kiểm tra các proxy.
// Trong Clean Architecture, đây là một interface ở tầng adapter
// giúp kết nối domain logic với implementation cụ thể
type BoKiemTra interface {
	KiemTraProxy(proxy model.Proxy, duongDanKiemTra string) (diaChi string, trangThai string, err error)
}

// BoKiemTraProxy là bộ điều hợp (adapter) để kiểm tra các proxy.
// Đóng vai trò trung gian giữa usecase và implementation cụ thể
type BoKiemTraProxy struct {
	boGhiNhatKy *utils.Logger // boGhiNhatKy: bộ ghi nhật ký
	boKiemTra   BoKiemTra     // boKiemTra: bộ kiểm tra - triển khai cụ thể (ví dụ: bộ kiểm tra IP)
}

// TaoBoKiemTraProxy tạo một bộ điều hợp (adapter) mới để kiểm tra proxy.
// Theo nguyên tắc Dependency Injection của Clean Architecture
func TaoBoKiemTraProxy(boGhiNhatKy *utils.Logger, boKiemTra BoKiemTra) *BoKiemTraProxy {
	return &BoKiemTraProxy{
		boGhiNhatKy: boGhiNhatKy,
		boKiemTra:   boKiemTra,
	}
}

// CheckProxy kiểm tra một proxy để xác định địa chỉ IP thực được sử dụng.
// Phương thức này ủy quyền kiểm tra cho triển khai cụ thể (boKiemTra)
func (c *BoKiemTraProxy) CheckProxy(proxy model.Proxy, duongDanKiemTra string) (diaChi string, trangThai string, err error) {
	diaChi, trangThai, err = c.boKiemTra.KiemTraProxy(proxy, duongDanKiemTra)
	if err != nil {
		c.boGhiNhatKy.Errorf("Kiểm tra proxy thất bại %s://%s:%s: %v", proxy.GiaoDien, proxy.DiaChi, proxy.Cong, err)
		return "", trangThai, err
	}
	c.boGhiNhatKy.Infof("Proxy %s://%s:%s trả về IP: %s", proxy.GiaoDien, proxy.DiaChi, proxy.Cong, diaChi)
	return diaChi, trangThai, nil
}
