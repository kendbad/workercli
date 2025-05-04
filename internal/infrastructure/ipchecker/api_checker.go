package ipchecker

import (
	"encoding/json"
	"fmt"
	"strings"
	"workercli/internal/adapter/httpclient"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

type APIChecker struct {
	client httpclient.HTTPClient
	logger *utils.Logger
}

func NewAPIChecker(clientType string, logger *utils.Logger) *APIChecker {
	client := httpclient.NewHTTPClient(clientType, logger)
	return &APIChecker{client: client, logger: logger}
}

func (c *APIChecker) CheckIP(proxy model.Proxy, checkURL string) (string, int, error) {
	// Đảm bảo checkURL là URL hợp lệ
	if !strings.HasPrefix(checkURL, "http://") && !strings.HasPrefix(checkURL, "https://") {
		checkURL = "http://" + checkURL
	}

	body, statusCode, err := c.client.DoRequest(proxy, checkURL)
	if err != nil {
		c.logger.Errorf("IP check failed for proxy %s://%s:%s: %v", proxy.Protocol, proxy.IP, proxy.Port, err)
		return "", statusCode, err
	}

	var ipResp struct {
		IP string `json:"ip"`
	}
	if err := json.Unmarshal(body, &ipResp); err != nil {
		c.logger.Errorf("Failed to decode JSON response: %v", err)
		return "", statusCode, fmt.Errorf("Không thể giải mã JSON: %v", err)
	}

	return ipResp.IP, statusCode, nil
}
