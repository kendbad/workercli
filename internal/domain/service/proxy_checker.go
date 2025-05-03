package service

import (
	"workercli/internal/domain/model"
)

type ProxyChecker interface {
	CheckProxy(proxy model.Proxy, checkURL string) (model.ProxyResult, error)
}
