package proxy

import (
	"bufio"
	"os"
	"strings"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

type Reader interface {
	ReadProxies(filePath string) ([]model.Proxy, error)
}

type ProxyReader struct {
	logger *utils.Logger
}

func NewProxyReader(logger *utils.Logger) *ProxyReader {
	return &ProxyReader{logger: logger}
}

func (r *ProxyReader) ReadProxies(filePath string) ([]model.Proxy, error) {
	filePath = utils.AutoPath(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		r.logger.Errorf("Không thể mở file proxy: %v", err)
		return nil, err
	}
	defer file.Close()

	var proxies []model.Proxy
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, "://")
		if len(parts) != 2 {
			r.logger.Warnf("Proxy không hợp lệ: %s", line)
			continue
		}
		protocol := parts[0]
		addrParts := strings.Split(parts[1], ":")
		if len(addrParts) != 2 {
			r.logger.Warnf("Proxy không hợp lệ: %s", line)
			continue
		}
		proxies = append(proxies, model.Proxy{
			Protocol: protocol,
			IP:       addrParts[0],
			Port:     addrParts[1],
		})
	}
	if err := scanner.Err(); err != nil {
		r.logger.Errorf("Lỗi đọc file proxy: %v", err)
		return nil, err
	}
	r.logger.Infof("Đọc được %d proxy từ %s", len(proxies), filePath)
	return proxies, nil
}
