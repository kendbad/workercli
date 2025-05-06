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

type KiemTraTrungGian struct {
	boDoc           *proxy.ProxyReader
	boKiemTra       *proxy.ProxyChecker
	duongDanKiemTra string
	nhomXuLy        *worker.NhomXuLy
	boGhiNhatKy     *utils.Logger
}

func TaoBoKiemTraTrungGian(boDoc *proxy.ProxyReader, boKiemTra *proxy.ProxyChecker, duongDanKiemTra string, soLuongXuLy int, boGhiNhatKy *utils.Logger) *KiemTraTrungGian {
	// Đảm bảo duongDanKiemTra có giao thức
	if !strings.HasPrefix(duongDanKiemTra, "http://") && !strings.HasPrefix(duongDanKiemTra, "https://") {
		duongDanKiemTra = "http://" + duongDanKiemTra
	}
	boXuLyTrungGian := &BoXuLyTacVuTrungGian{boKiemTra, duongDanKiemTra, boGhiNhatKy}
	return &KiemTraTrungGian{
		boDoc:           boDoc,
		boKiemTra:       boKiemTra,
		duongDanKiemTra: duongDanKiemTra,
		nhomXuLy:        worker.TaoNhomXuLy(soLuongXuLy, boXuLyTrungGian, boGhiNhatKy),
		boGhiNhatKy:     boGhiNhatKy,
	}
}

type BoXuLyTacVuTrungGian struct {
	boKiemTra       *proxy.ProxyChecker
	duongDanKiemTra string
	boGhiNhatKy     *utils.Logger
}

func (p *BoXuLyTacVuTrungGian) XuLyTacVu(tacVu model.TacVu) (model.KetQua, error) {
	trungGian, err := proxi.ParseProxy(tacVu.MaTacVu)
	if err != nil {
		p.boGhiNhatKy.Errorf("Định dạng proxy không hợp lệ: %s", tacVu.MaTacVu)
		return model.KetQua{MaTacVu: tacVu.MaTacVu, TrangThai: "Thất bại", LoiXayRa: err.Error()}, err
	}

	diaChi, trangThai, err := p.boKiemTra.CheckProxy(trungGian, p.duongDanKiemTra)
	if err != nil {
		p.boGhiNhatKy.Errorf("Kiểm tra proxy thất bại %s: %v", tacVu.MaTacVu, err)
		return model.KetQua{MaTacVu: tacVu.MaTacVu, TrangThai: trangThai, LoiXayRa: err.Error()}, err
	}

	p.boGhiNhatKy.Infof("Proxy %s trả về IP: %s", tacVu.MaTacVu, diaChi)
	return model.KetQua{MaTacVu: tacVu.MaTacVu, TrangThai: trangThai}, nil
}

func (uc *KiemTraTrungGian) ThucThi(duongDanFileTrungGian string) ([]model.KetQuaTrungGian, error) {
	danhSachTrungGian, err := uc.boDoc.ReadProxies(duongDanFileTrungGian)
	if err != nil {
		uc.boGhiNhatKy.Errorf("Không đọc được danh sách proxy: %v", err)
		return nil, err
	}

	uc.nhomXuLy.BatDau()
	ketQua := make([]model.KetQuaTrungGian, 0, len(danhSachTrungGian))
	kenhKetQuaTrungGian := make(chan model.KetQuaTrungGian, len(danhSachTrungGian))

	for _, p := range danhSachTrungGian {
		tacVu := model.TacVu{
			MaTacVu: fmt.Sprintf("%s://%s:%s", p.GiaoDien, p.DiaChi, p.Cong),
			DuLieu:  p.GiaoDien,
		}
		uc.nhomXuLy.NopTacVu(tacVu)

		go func(trungGian model.TrungGian) {
			ketQuaTacVu := <-uc.nhomXuLy.KetQua()
			ketQuaDon := model.KetQuaTrungGian{TrungGian: trungGian, TrangThai: ketQuaTacVu.TrangThai, LoiXayRa: ketQuaTacVu.LoiXayRa}
			if ketQuaTacVu.TrangThai == "Thành công" {
				if diaChi, _, err := uc.boKiemTra.CheckProxy(trungGian, uc.duongDanKiemTra); err == nil {
					ketQuaDon.DiaChi = diaChi
				} else {
					ketQuaDon.DiaChi = "Lấy IP thất bại"
				}
			} else {
				ketQuaDon.DiaChi = "Gửi tác vụ thất bại"
			}
			kenhKetQuaTrungGian <- ketQuaDon
		}(p)
	}

	for i := 0; i < len(danhSachTrungGian); i++ {
		ketQua = append(ketQua, <-kenhKetQuaTrungGian)
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
			r.TrungGian.GiaoDien, r.TrungGian.DiaChi, r.TrungGian.Cong, r.DiaChi, trangThai)
	}

	uc.boGhiNhatKy.Infof("Kết quả đã được lưu vào %s", duongDanFileKetQua)
	return ketQua, nil
}
