package tview

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

func (vm *ViewModel) GetResults() []model.KetQua {
	return vm.KetQua
}

func (vm *ViewModel) GetProxies() []model.Proxy {
	return vm.DanhSachProxy
}

func (vm *ViewModel) GetResult(index int) model.KetQua {
	return vm.KetQua[index]
}

func (vm *ViewModel) GetProxy(index int) model.Proxy {
	return vm.DanhSachProxy[index]
}

func (vm *ViewModel) GetResultCount() int {
	return len(vm.KetQua)
}

func (vm *ViewModel) GetProxyCount() int {
	return len(vm.DanhSachProxy)
}
