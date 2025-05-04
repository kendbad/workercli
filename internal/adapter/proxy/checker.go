package proxy

import (
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// Checker defines the interface for checking proxies.
type Checker interface {
	CheckProxy(proxy model.Proxy, checkURL string) (ip string, status string, err error)
}

// ProxyChecker is the adapter for checking proxies.
type ProxyChecker struct {
	logger  *utils.Logger
	checker Checker // Specific implementation (e.g., IP checker)
}

func NewProxyChecker(logger *utils.Logger, checker Checker) *ProxyChecker {
	return &ProxyChecker{
		logger:  logger,
		checker: checker,
	}
}

func (c *ProxyChecker) CheckProxy(proxy model.Proxy, checkURL string) (ip string, status string, err error) {
	ip, status, err = c.checker.CheckProxy(proxy, checkURL)
	if err != nil {
		c.logger.Errorf("Proxy check failed %s://%s:%s: %v", proxy.Protocol, proxy.IP, proxy.Port, err)
		return "", status, err
	}
	c.logger.Infof("Proxy %s://%s:%s returned IP: %s", proxy.Protocol, proxy.IP, proxy.Port, ip)
	return ip, status, nil
}
