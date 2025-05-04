package proxy

import (
	"fmt"
	"workercli/internal/adapter/ipchecker"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

type IPChecker struct {
	logger    *utils.Logger
	ipChecker ipchecker.IPChecker
}

func NewIPChecker(logger *utils.Logger, clientType string) *IPChecker {
	ipChecker := ipchecker.NewIPChecker(clientType, logger)
	return &IPChecker{logger: logger, ipChecker: ipChecker}
}

func (c *IPChecker) CheckProxy(proxy model.Proxy, checkURL string) (ip string, status string, err error) {
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

	return ip, "Success", nil
}
