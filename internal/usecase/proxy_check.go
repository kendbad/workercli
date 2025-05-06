package usecase

import (
	"fmt"
	"os"
	"strings"
	"workercli/internal/adapter/proxy"
	"workercli/internal/adapter/worker"
	"workercli/internal/domain/model"
	proxi "workercli/internal/infrastructure/proxy"
	"workercli/pkg/utils"
)

// ParseProxyFunc là biến hàm cho phép inject dependency và mock trong test
// Clean Architecture: Dependency Rule - Tầng ngoài phụ thuộc vào tầng trong
var ParseProxyFunc = proxi.ParseProxy

// KiemTraProxy là usecase để kiểm tra danh sách proxy.
// Trong Clean Architecture, usecase chứa logic nghiệp vụ cụ thể
// và điều phối các thành phần khác để thực hiện một trường hợp sử dụng.
// Usecase nằm ở tầng trong, không phụ thuộc vào infrastructure hay adapter
type KiemTraProxy struct {
	BoDoc           *proxy.BoDocProxy     // boDoc: bộ đọc proxy
	BoKiemTra       *proxy.BoKiemTraProxy // boKiemTra: bộ kiểm tra proxy
	DuongDanKiemTra string                // duongDanKiemTra: đường dẫn kiểm tra
	NhomXuLy        *worker.NhomXuLy      // nhomXuLy: nhóm xử lý - quản lý pool worker
	BoGhiNhatKy     *utils.Logger         // boGhiNhatKy: bộ ghi nhật ký
}

// TaoBoKiemTraProxy tạo một usecase mới để kiểm tra proxy.
// Áp dụng Dependency Injection để tiêm phụ thuộc qua tham số
// Clean Architecture: Factory Pattern - Tạo usecase với các dependency đã chuẩn bị
func TaoBoKiemTraProxy(boDoc *proxy.BoDocProxy, boKiemTra *proxy.BoKiemTraProxy, duongDanKiemTra string, soLuongXuLy int, boGhiNhatKy *utils.Logger) *KiemTraProxy {
	// Đảm bảo duongDanKiemTra có giao thức
	if !strings.HasPrefix(duongDanKiemTra, "http://") && !strings.HasPrefix(duongDanKiemTra, "https://") {
		duongDanKiemTra = "http://" + duongDanKiemTra
	}
	boXuLyProxy := &BoXuLyTacVuProxy{boKiemTra, duongDanKiemTra, boGhiNhatKy}
	return &KiemTraProxy{
		BoDoc:           boDoc,
		BoKiemTra:       boKiemTra,
		DuongDanKiemTra: duongDanKiemTra,
		NhomXuLy:        worker.TaoNhomXuLy(soLuongXuLy, boXuLyProxy, boGhiNhatKy),
		BoGhiNhatKy:     boGhiNhatKy,
	}
}

// BoXuLyTacVuProxy là một adapter nội bộ để xử lý tác vụ kiểm tra proxy.
// Đây là một ví dụ về adapter pattern trong Clean Architecture,
// cho phép usecase tương tác với các worker thông qua một giao diện chung.
// Clean Architecture: Port & Adapter - BoXuLyTacVuProxy là adapter cho worker port
type BoXuLyTacVuProxy struct {
	boKiemTra       *proxy.BoKiemTraProxy // boKiemTra: bộ kiểm tra proxy
	duongDanKiemTra string                // duongDanKiemTra: đường dẫn kiểm tra
	boGhiNhatKy     *utils.Logger         // boGhiNhatKy: bộ ghi nhật ký
}

// XuLyTacVu xử lý một tác vụ kiểm tra proxy.
// Implements giao diện worker.TaskProcessor để worker có thể gọi
// Clean Architecture: Interface Segregation - Chỉ expose các phương thức cần thiết
func (p *BoXuLyTacVuProxy) XuLyTacVu(tacVu model.TacVu) (model.KetQua, error) {
	proxy, err := ParseProxyFunc(tacVu.MaTacVu)
	if err != nil {
		p.boGhiNhatKy.Errorf("Định dạng proxy không hợp lệ: %s", tacVu.MaTacVu)
		return model.KetQua{MaTacVu: tacVu.MaTacVu, TrangThai: "Thất bại", LoiXayRa: err.Error()}, err
	}

	diaChi, trangThai, err := p.boKiemTra.CheckProxy(proxy, p.duongDanKiemTra)
	if err != nil {
		p.boGhiNhatKy.Errorf("Kiểm tra proxy thất bại %s: %v", tacVu.MaTacVu, err)
		return model.KetQua{MaTacVu: tacVu.MaTacVu, TrangThai: trangThai, LoiXayRa: err.Error()}, err
	}

	p.boGhiNhatKy.Infof("Proxy %s trả về IP: %s", tacVu.MaTacVu, diaChi)
	return model.KetQua{MaTacVu: tacVu.MaTacVu, TrangThai: trangThai}, nil
}

// ThucThi thực hiện việc kiểm tra tất cả proxy từ một tệp tin.
// Đây là phương thức chính của usecase, điều phối toàn bộ luồng xử lý
// Clean Architecture: Use Case Interactor - Chứa business logic của ứng dụng
func (uc *KiemTraProxy) ThucThi(duongDanFileProxy string) ([]model.KetQuaProxy, error) {
	danhSachProxy, err := uc.BoDoc.ReadProxies(duongDanFileProxy)
	if err != nil {
		uc.BoGhiNhatKy.Errorf("Không đọc được danh sách proxy: %v", err)
		return nil, err
	}

	uc.NhomXuLy.BatDau()
	ketQua := make([]model.KetQuaProxy, 0, len(danhSachProxy))
	kenhKetQuaProxy := make(chan model.KetQuaProxy, len(danhSachProxy))

	for _, p := range danhSachProxy {
		tacVu := model.TacVu{
			MaTacVu: fmt.Sprintf("%s://%s:%s", p.GiaoDien, p.DiaChi, p.Cong),
			DuLieu:  p.GiaoDien,
		}
		uc.NhomXuLy.NopTacVu(tacVu)

		go func(proxy model.Proxy) {
			ketQuaTacVu := <-uc.NhomXuLy.KetQua()
			ketQuaDon := model.KetQuaProxy{Proxy: proxy, TrangThai: ketQuaTacVu.TrangThai, LoiXayRa: ketQuaTacVu.LoiXayRa}
			if ketQuaTacVu.TrangThai == "Thành công" {
				if diaChi, _, err := uc.BoKiemTra.CheckProxy(proxy, uc.DuongDanKiemTra); err == nil {
					ketQuaDon.DiaChi = diaChi
				} else {
					ketQuaDon.DiaChi = "Lấy IP thất bại"
				}
			} else {
				ketQuaDon.DiaChi = "Gửi tác vụ thất bại"
			}
			kenhKetQuaProxy <- ketQuaDon
		}(p)
	}

	for i := 0; i < len(danhSachProxy); i++ {
		ketQua = append(ketQua, <-kenhKetQuaProxy)
	}

	uc.NhomXuLy.Dung()

	duongDanFileKetQua := utils.AutoPath("output/ket_qua_proxy.txt")
	tep, err := os.Create(duongDanFileKetQua)
	if err != nil {
		uc.BoGhiNhatKy.Errorf("Không tạo được file kết quả: %v", err)
		return ketQua, err
	}
	defer tep.Close()

	for _, r := range ketQua {
		trangThai := r.TrangThai
		if r.LoiXayRa != "" {
			trangThai += " (" + r.LoiXayRa + ")"
		}
		fmt.Fprintf(tep, "Proxy: %s://%s:%s, IP: %s, Trạng thái: %s\n",
			r.Proxy.GiaoDien, r.Proxy.DiaChi, r.Proxy.Cong, r.DiaChi, trangThai)
	}

	uc.BoGhiNhatKy.Infof("Kết quả đã được lưu vào %s", duongDanFileKetQua)
	return ketQua, nil
}
