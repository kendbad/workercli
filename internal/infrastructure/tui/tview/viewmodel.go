package tview

import (
	"workercli/internal/domain/model"
)

type ViewModel struct {
	Results []model.Result
	Proxies []model.Proxy
}

func (vm *ViewModel) UpdateResults(results []model.Result) {
	vm.Results = results
}

func (vm *ViewModel) UpdateProxies(proxies []model.Proxy) {
	vm.Proxies = proxies
}

func (vm *ViewModel) GetResults() []model.Result {
	return vm.Results
}

func (vm *ViewModel) GetProxies() []model.Proxy {
	return vm.Proxies
}

func (vm *ViewModel) GetResult(index int) model.Result {
	return vm.Results[index]
}

func (vm *ViewModel) GetProxy(index int) model.Proxy {
	return vm.Proxies[index]
}

func (vm *ViewModel) GetResultCount() int {
	return len(vm.Results)
}

func (vm *ViewModel) GetProxyCount() int {
	return len(vm.Proxies)
}
