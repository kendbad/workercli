package proxy

import (
	"bufio"
	"os"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

type FileReader struct {
	logger *utils.Logger
}

func NewFileReader(logger *utils.Logger) *FileReader {
	return &FileReader{logger: logger}
}

func (r *FileReader) ReadProxies(filePath string) ([]model.Proxy, error) {
	filePath = utils.AutoPath(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		r.logger.Errorf("Không thể mở file proxy %s: %v", filePath, err)
		return nil, err
	}
	defer file.Close()

	var proxies []model.Proxy
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		proxy, err := ParseProxy(line)
		if err != nil {
			r.logger.Warnf("Bỏ qua proxy không hợp lệ %s: %v", line, err)
			continue
		}
		proxies = append(proxies, proxy)
	}

	if err := scanner.Err(); err != nil {
		r.logger.Errorf("Lỗi khi đọc file proxy %s: %v", filePath, err)
		return nil, err
	}

	return proxies, nil
}
