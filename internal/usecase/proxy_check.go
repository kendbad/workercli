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

type KiemTraProxy struct {
	boDoc           *proxy.ProxyReader
	boKiemTra       *proxy.ProxyChecker
	duongDanKiemTra string
	nhomXuLy        *worker.NhomXuLy
	boGhiNhatKy     *utils.Logger
}

func TaoBoKiemTraProxy(boDoc *proxy.ProxyReader, boKiemTra *proxy.ProxyChecker, duongDanKiemTra string, soLuongXuLy int, boGhiNhatKy *utils.Logger) *KiemTraProxy {
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

type BoXuLyTacVuProxy struct {
	boKiemTra       *proxy.ProxyChecker
	duongDanKiemTra string
	boGhiNhatKy     *utils.Logger
}

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
