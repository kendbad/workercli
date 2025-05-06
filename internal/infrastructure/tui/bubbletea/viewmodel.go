package bubbletea

import (
	"workercli/internal/domain/model"
)

type ViewModel struct {
	KetQua        []model.KetQua
	DanhSachProxy []model.Proxy
}

func (vm *ViewModel) UpdateResults(ketQua []model.KetQua) {
	vm.KetQua = ketQua
}

func (vm *ViewModel) UpdateProxies(danhSachProxy []model.Proxy) {
	vm.DanhSachProxy = danhSachProxy
}
