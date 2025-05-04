package ipchecker

import (
	"encoding/json"
	"fmt"
	"workercli/internal/domain/model"
	"workercli/internal/httpclient"
	"workercli/pkg/utils"
)

type IPChecker struct {
	client httpclient.HTTPClient
	logger *utils.Logger
}

func NewIPChecker(clientType string, logger *utils.Logger) *IPChecker {
	var client httpclient.HTTPClient
	switch clientType {
	case "fasthttp":
		client = httpclient.NewFastHTTPClient(logger)
	case "nethttp":
		client = httpclient.NewNetHTTPClient(logger)
	default:
		logger.Warnf("Unknown client type %s, falling back to fasthttp", clientType)
		client = httpclient.NewFastHTTPClient(logger)
	}
	return &IPChecker{client: client, logger: logger}
}

func (c *IPChecker) CheckIP(proxy model.Proxy, checkURL string) (string, int, error) {
	body, statusCode, err := c.client.DoRequest(proxy, checkURL)
	if err != nil {
		return "", statusCode, err
	}

	var ipResp struct {
		IP string `json:"ip"`
	}
	if err := json.Unmarshal(body, &ipResp); err != nil {
		return "", statusCode, fmt.Errorf("Không thể giải mã JSON: %v", err)
	}

	return ipResp.IP, statusCode, nil
}
