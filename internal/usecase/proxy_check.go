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

// KiemTraProxy là usecase để kiểm tra danh sách proxy.
// Trong Clean Architecture, usecase chứa logic nghiệp vụ cụ thể
// và điều phối các thành phần khác để thực hiện một trường hợp sử dụng.
type KiemTraProxy struct {
	boDoc           *proxy.ProxyReader    // boDoc: bộ đọc proxy
	boKiemTra       *proxy.BoKiemTraProxy // boKiemTra: bộ kiểm tra proxy
	duongDanKiemTra string                // duongDanKiemTra: đường dẫn kiểm tra
	nhomXuLy        *worker.NhomXuLy      // nhomXuLy: nhóm xử lý - quản lý pool worker
	boGhiNhatKy     *utils.Logger         // boGhiNhatKy: bộ ghi nhật ký
}

// TaoBoKiemTraProxy tạo một usecase mới để kiểm tra proxy.
// Áp dụng Dependency Injection để tiêm phụ thuộc qua tham số
func TaoBoKiemTraProxy(boDoc *proxy.ProxyReader, boKiemTra *proxy.BoKiemTraProxy, duongDanKiemTra string, soLuongXuLy int, boGhiNhatKy *utils.Logger) *KiemTraProxy {
	// Đảm bảo duongDanKiemTra có giao thức
	if !strings.HasPrefix(duongDanKiemTra, "http://") && !strings.HasPrefix(duongDanKiemTra, "https://") {
		duongDanKiemTra = "http://" + duongDanKiemTra
	}
	boXuLyProxy := &BoXuLyTacVuProxy{boKiemTra, duongDanKiemTra, boGhiNhatKy}
	return &KiemTraProxy{
		boDoc:           boDoc,
		boKiemTra:       boKiemTra,
		duongDanKiemTra: duongDanKiemTra,
		nhomXuLy:        worker.TaoNhomXuLy(soLuongXuLy, boXuLyProxy, boGhiNhatKy),
		boGhiNhatKy:     boGhiNhatKy,
	}
}

// BoXuLyTacVuProxy là một adapter nội bộ để xử lý tác vụ kiểm tra proxy.
// Đây là một ví dụ về adapter pattern trong Clean Architecture,
// cho phép usecase tương tác với các worker thông qua một giao diện chung.
type BoXuLyTacVuProxy struct {
	boKiemTra       *proxy.BoKiemTraProxy // boKiemTra: bộ kiểm tra proxy
	duongDanKiemTra string                // duongDanKiemTra: đường dẫn kiểm tra
	boGhiNhatKy     *utils.Logger         // boGhiNhatKy: bộ ghi nhật ký
}

// XuLyTacVu xử lý một tác vụ kiểm tra proxy.
// Implements giao diện worker.TaskProcessor để worker có thể gọi
func (p *BoXuLyTacVuProxy) XuLyTacVu(tacVu model.TacVu) (model.KetQua, error) {
	proxy, err := proxi.ParseProxy(tacVu.MaTacVu)
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
func (uc *KiemTraProxy) ThucThi(duongDanFileProxy string) ([]model.KetQuaProxy, error) {
	danhSachProxy, err := uc.boDoc.ReadProxies(duongDanFileProxy)
	if err != nil {
		uc.boGhiNhatKy.Errorf("Không đọc được danh sách proxy: %v", err)
		return nil, err
	}

	uc.nhomXuLy.BatDau()
	ketQua := make([]model.KetQuaProxy, 0, len(danhSachProxy))
	kenhKetQuaProxy := make(chan model.KetQuaProxy, len(danhSachProxy))

	for _, p := range danhSachProxy {
		tacVu := model.TacVu{
			MaTacVu: fmt.Sprintf("%s://%s:%s", p.GiaoDien, p.DiaChi, p.Cong),
			DuLieu:  p.GiaoDien,
		}
		uc.nhomXuLy.NopTacVu(tacVu)

		go func(proxy model.Proxy) {
			ketQuaTacVu := <-uc.nhomXuLy.KetQua()
			ketQuaDon := model.KetQuaProxy{Proxy: proxy, TrangThai: ketQuaTacVu.TrangThai, LoiXayRa: ketQuaTacVu.LoiXayRa}
			if ketQuaTacVu.TrangThai == "Thành công" {
				if diaChi, _, err := uc.boKiemTra.CheckProxy(proxy, uc.duongDanKiemTra); err == nil {
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

	uc.nhomXuLy.Dung()

	duongDanFileKetQua := utils.AutoPath("output/ket_qua_proxy.txt")
	tep, err := os.Create(duongDanFileKetQua)
	if err != nil {
		uc.boGhiNhatKy.Errorf("Không tạo được file kết quả: %v", err)
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

	uc.boGhiNhatKy.Infof("Kết quả đã được lưu vào %s", duongDanFileKetQua)
	return ketQua, nil
}
