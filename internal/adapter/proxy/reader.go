package proxy

import (
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// Reader defines the interface for reading proxies.
type Reader interface {
	ReadProxies(source string) ([]model.Proxy, error)
}

// ProxyReader is the adapter for reading proxies.
type ProxyReader struct {
	logger *utils.Logger
	reader Reader // Specific implementation (e.g., file reader)
}

func NewProxyReader(logger *utils.Logger, reader Reader) *ProxyReader {
	return &ProxyReader{
		logger: logger,
		reader: reader,
	}
}

func (r *ProxyReader) ReadProxies(source string) ([]model.Proxy, error) {
	proxies, err := r.reader.ReadProxies(source)
	if err != nil {
		r.logger.Errorf("Failed to read proxies from %s: %v", source, err)
		return nil, err
	}
	r.logger.Infof("Read %d proxies from %s", len(proxies), source)
	return proxies, nil
}
