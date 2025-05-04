package proxy

import (
	"fmt"
	"workercli/internal/domain/model"
	"workercli/internal/ipchecker"
	"workercli/pkg/utils"
)

type Checker struct {
	logger    *utils.Logger
	ipChecker *ipchecker.IPChecker
}

func NewChecker(logger *utils.Logger, clientType string) *Checker {
	ipChecker := ipchecker.NewIPChecker(clientType, logger)
	return &Checker{logger: logger, ipChecker: ipChecker}
}

func (c *Checker) CheckProxy(proxy model.Proxy, checkURL string) (ip string, status string, err error) {
	ip, statusCode, err := c.ipChecker.CheckIP(proxy, checkURL)
	if err != nil {
		c.logger.Errorf("Proxy %s://%s:%s failed: %v", proxy.Protocol, proxy.IP, proxy.Port, err)
		return "", fmt.Sprintf("Failed (%v)", err), err
	}

	if statusCode != 200 {
		err = fmt.Errorf("mã trạng thái: %d", statusCode)
		c.logger.Errorf("Proxy %s://%s:%s returned status: %d", proxy.Protocol, proxy.IP, proxy.Port, statusCode)
		return "", fmt.Sprintf("Failed (status: %d)", statusCode), err
	}

	c.logger.Infof("Proxy %s://%s:%s returned IP: %s", proxy.Protocol, proxy.IP, proxy.Port, ip)
	return ip, "Success", nil
}
