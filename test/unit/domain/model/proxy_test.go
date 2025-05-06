package model_test

import (
	"testing"
	"workercli/internal/domain/model"
)

func TestProxy(t *testing.T) {
	// Khởi tạo struct Proxy để kiểm tra
	proxy := model.Proxy{
		GiaoDien: "http",
		DiaChi:   "192.168.1.1",
		Cong:     "8080",
	}

	// Kiểm tra các trường dữ liệu
	if proxy.GiaoDien != "http" {
		t.Errorf("Proxy.GiaoDien = %s; muốn có 'http'", proxy.GiaoDien)
	}

	if proxy.DiaChi != "192.168.1.1" {
		t.Errorf("Proxy.DiaChi = %s; muốn có '192.168.1.1'", proxy.DiaChi)
	}

	if proxy.Cong != "8080" {
		t.Errorf("Proxy.Cong = %s; muốn có '8080'", proxy.Cong)
	}
}

func TestKetQuaProxy(t *testing.T) {
	// Khởi tạo Proxy
	proxy := model.Proxy{
		GiaoDien: "http",
		DiaChi:   "192.168.1.1",
		Cong:     "8080",
	}

	// Khởi tạo KetQuaProxy
	ketQua := model.KetQuaProxy{
		Proxy:     proxy,
		DiaChi:    "123.45.67.89",
		TrangThai: "Thành công",
		LoiXayRa:  "",
	}

	// Kiểm tra các trường dữ liệu
	if ketQua.DiaChi != "123.45.67.89" {
		t.Errorf("KetQuaProxy.DiaChi = %s; muốn có '123.45.67.89'", ketQua.DiaChi)
	}

	if ketQua.TrangThai != "Thành công" {
		t.Errorf("KetQuaProxy.TrangThai = %s; muốn có 'Thành công'", ketQua.TrangThai)
	}

	// Kiểm tra thông tin Proxy trong KetQuaProxy
	if ketQua.Proxy.GiaoDien != "http" {
		t.Errorf("KetQuaProxy.Proxy.GiaoDien = %s; muốn có 'http'", ketQua.Proxy.GiaoDien)
	}
}
