package worker

import (
	"sync"
	"workercli/internal/domain/model"
	"workercli/internal/domain/service"
	"workercli/pkg/utils"
)

// NhomXuLy quản lý một tập hợp các người xử lý
type NhomXuLy struct {
	danhSachTacVu  chan model.TacVu
	danhSachKetQua chan model.KetQua
	boXuLy         service.BoXuLyTacVu
	soLuongXuLy    int
	boGhiNhatKy    *utils.Logger
	kenhDung       chan struct{}
	nhomCho        sync.WaitGroup
}

// TaoNhomXuLy tạo một nhóm xử lý mới
func TaoNhomXuLy(soLuongXuLy int, boXuLy service.BoXuLyTacVu, boGhiNhatKy *utils.Logger) *NhomXuLy {
	return &NhomXuLy{
		danhSachTacVu:  make(chan model.TacVu, 1000), // Hàng đợi lớn cho số lượng tác vụ lớn
		danhSachKetQua: make(chan model.KetQua, 1000),
		boXuLy:         boXuLy,
		soLuongXuLy:    soLuongXuLy,
		boGhiNhatKy:    boGhiNhatKy,
		kenhDung:       make(chan struct{}),
	}
}

// BatDau khởi động tất cả bộ xử lý
func (p *NhomXuLy) BatDau() {
	p.boGhiNhatKy.Infof("Khởi động nhóm với %d bộ xử lý", p.soLuongXuLy)
	p.nhomCho.Add(p.soLuongXuLy)
	for i := 0; i < p.soLuongXuLy; i++ {
		nguoiXuLy := TaoNguoiXuLy(i+1, p.danhSachTacVu, p.danhSachKetQua, p.boXuLy, p.boGhiNhatKy)
		go func() {
			defer p.nhomCho.Done()
			nguoiXuLy.Chay(p.kenhDung)
		}()
	}
}

// NopTacVu gửi một tác vụ vào hàng đợi
func (p *NhomXuLy) NopTacVu(tacVu model.TacVu) {
	p.danhSachTacVu <- tacVu
}

// KetQua trả về kênh kết quả
func (p *NhomXuLy) KetQua() <-chan model.KetQua {
	return p.danhSachKetQua
}

// Dung dừng tất cả bộ xử lý
func (p *NhomXuLy) Dung() {
	p.boGhiNhatKy.Info("Dừng nhóm xử lý")
	close(p.kenhDung)
	close(p.danhSachTacVu)
	p.nhomCho.Wait()
	close(p.danhSachKetQua)
}
