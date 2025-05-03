package bubbletea

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
