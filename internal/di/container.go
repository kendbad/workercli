package di

import (
	"sync"

	"workercli/internal/adapter/input"
	proxyiface "workercli/internal/adapter/proxy"
	tuiadapter "workercli/internal/adapter/tui"
	"workercli/internal/config"
	"workercli/internal/domain/model"
	"workercli/internal/domain/service"
	"workercli/internal/infrastructure/proxy"
	"workercli/internal/infrastructure/task"
	tuiinfra "workercli/internal/infrastructure/tui"
	"workercli/internal/usecase"
	"workercli/pkg/utils"
)

// Container quản lý và cung cấp tất cả các dependencies cho ứng dụng
type Container struct {
	cauHinh     *config.Config
	boGhiNhatKy *utils.Logger
	kenhDungLai chan struct{}
	boGiaoDien  *tuiinfra.TUIUseCase

	// Các thành phần cho proxy
	boDocTep       *proxy.FileReader
	boDocProxy     *proxyiface.BoDocProxy
	boKiemTraIP    *proxy.BoKiemTraIP
	boKiemTraProxy *proxyiface.BoKiemTraProxy
	boKiemTra      *usecase.KiemTraProxy

	// Các thành phần cho task
	boXuLy         service.BoXuLyTacVu
	boDocDauVao    input.Reader
	boXuLyHangLoat *usecase.XuLyHangLoatTacVu

	// Cờ kiểm tra khởi tạo
	daKhoiTao bool
	mutex     sync.Mutex
}

// NewContainer tạo container DI mới
func NewContainer() *Container {
	return &Container{
		kenhDungLai: make(chan struct{}),
		daKhoiTao:   false,
	}
}

// KhoiTao khởi tạo container với cấu hình
func (c *Container) KhoiTao(thuMucCauHinh string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.daKhoiTao {
		return nil
	}

	// Tải cấu hình
	cauHinh, err := config.Load(thuMucCauHinh)
	if err != nil {
		return err
	}
	c.cauHinh = cauHinh

	// Khởi tạo bộ ghi nhật ký
	boGhiNhatKy, err := utils.NewLogger("configs/logger.yaml")
	if err != nil {
		return err
	}
	c.boGhiNhatKy = boGhiNhatKy

	c.daKhoiTao = true
	return nil
}

// DungContainer đóng container và giải phóng tài nguyên
func (c *Container) DungContainer() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.boGiaoDien != nil {
		c.boGiaoDien.Close()
	}

	close(c.kenhDungLai)
	c.daKhoiTao = false
}

// LayBoGhiNhatKy trả về đối tượng logger
func (c *Container) LayBoGhiNhatKy() *utils.Logger {
	return c.boGhiNhatKy
}

// LayCauHinh trả về cấu hình
func (c *Container) LayCauHinh() *config.Config {
	return c.cauHinh
}

// LayKenhDungLai trả về kênh dừng
func (c *Container) LayKenhDungLai() chan struct{} {
	return c.kenhDungLai
}

// LayBoGiaoDien trả về đối tượng giao diện người dùng
func (c *Container) LayBoGiaoDien() *tuiinfra.TUIUseCase {
	return c.boGiaoDien
}

// ThietLapGiaoDien thiết lập giao diện người dùng
func (c *Container) ThietLapGiaoDien(kieuGiaoDien string, loaiHienThi model.LoaiHienThi) error {
	// Khởi tạo nhà xưởng giao diện
	nhaXuongGiaoDien := tuiinfra.NewRendererFactory(c.boGhiNhatKy, kieuGiaoDien)
	var renderer tuiadapter.Renderer

	// Tạo bộ hiển thị phù hợp
	if loaiHienThi == model.LoaiKiemTraProxy {
		renderer = nhaXuongGiaoDien.CreateProxyRenderer(
			c.boGhiNhatKy,
			&[]model.KetQuaProxy{},
			&sync.Mutex{},
			make(chan model.KetQuaProxy, 100),
			c.kenhDungLai,
		)
	} else if loaiHienThi == model.LoaiXuLyTacVu {
		renderer = nhaXuongGiaoDien.CreateTaskRenderer(
			c.boGhiNhatKy,
			&[]model.KetQua{},
			&sync.Mutex{},
			make(chan model.KetQua, 100),
			c.kenhDungLai,
		)
	} else {
		return nil
	}

	// Khởi tạo TUI
	c.boGiaoDien = tuiinfra.NewTUIUseCase(c.boGhiNhatKy, kieuGiaoDien, renderer)
	return nil
}

// ThietLapKiemTraProxy thiết lập các thành phần cho việc kiểm tra proxy
func (c *Container) ThietLapKiemTraProxy(loaiKetNoi, kieuGiaoDien string) (*usecase.KiemTraProxy, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Khởi tạo các thành phần nếu chưa có
	if c.boDocTep == nil {
		c.boDocTep = proxy.NewFileReader(c.boGhiNhatKy)
	}

	if c.boDocProxy == nil {
		c.boDocProxy = proxyiface.NewProxyReader(c.boGhiNhatKy, c.boDocTep)
	}

	if c.boKiemTraIP == nil {
		c.boKiemTraIP = proxy.NewIPChecker(c.boGhiNhatKy, loaiKetNoi)
	}

	if c.boKiemTraProxy == nil {
		c.boKiemTraProxy = proxyiface.TaoBoKiemTraProxy(c.boGhiNhatKy, c.boKiemTraIP)
	}

	// Khởi tạo usecase
	if c.boKiemTra == nil {
		c.boKiemTra = usecase.TaoBoKiemTraProxy(
			c.boDocProxy,
			c.boKiemTraProxy,
			c.cauHinh.Proxy.CheckURL,
			c.cauHinh.Worker.Workers,
			c.boGhiNhatKy,
		)
	}

	// Thiết lập TUI nếu cần
	if kieuGiaoDien != "" && c.boGiaoDien == nil {
		if err := c.ThietLapGiaoDien(kieuGiaoDien, model.LoaiKiemTraProxy); err != nil {
			return nil, err
		}
	}

	return c.boKiemTra, nil
}

// ThietLapXuLyTacVu thiết lập các thành phần cho việc xử lý tác vụ
func (c *Container) ThietLapXuLyTacVu(kieuGiaoDien string) (*usecase.XuLyHangLoatTacVu, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Khởi tạo các thành phần nếu chưa có
	if c.boXuLy == nil {
		c.boXuLy = task.NewProcessor(c.boGhiNhatKy)
	}

	if c.boDocDauVao == nil {
		c.boDocDauVao = input.NewFileReader(c.boGhiNhatKy)
	}

	// Khởi tạo usecase
	if c.boXuLyHangLoat == nil {
		c.boXuLyHangLoat = usecase.TaoBoXuLyHangLoatTacVu(
			c.boDocDauVao,
			c.boXuLy,
			c.cauHinh.Worker.Workers,
			c.boGhiNhatKy,
		)
	}

	// Thiết lập TUI nếu cần
	if kieuGiaoDien != "" && c.boGiaoDien == nil {
		if err := c.ThietLapGiaoDien(kieuGiaoDien, model.LoaiXuLyTacVu); err != nil {
			return nil, err
		}
	}

	return c.boXuLyHangLoat, nil
}
