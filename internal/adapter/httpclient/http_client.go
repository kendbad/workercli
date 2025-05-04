package httpclient

import (
	"workercli/internal/domain/model"
	httpclient "workercli/internal/infrastructure/httpclient"
	"workercli/pkg/utils"
)

// HTTPClient defines the interface for sending HTTP requests.
type HTTPClient interface {
	DoRequest(proxy model.Proxy, url string) ([]byte, int, error)
}

// Factory creates an HTTPClient based on the client type.
func NewHTTPClient(clientType string, logger *utils.Logger) HTTPClient {
	switch clientType {
	case "fasthttp":
		return httpclient.NewFastHTTPClient(logger)
	case "nethttp":
		return httpclient.NewNetHTTPClient(logger)
	default:
		logger.Warnf("Unknown client type %s, falling back to fasthttp", clientType)
		return httpclient.NewFastHTTPClient(logger)
	}
}
