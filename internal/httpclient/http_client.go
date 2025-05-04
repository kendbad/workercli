package httpclient

import (
	"workercli/internal/domain/model"
)

// HTTPClient defines a common interface for making HTTP requests
type HTTPClient interface {
	DoRequest(proxy model.Proxy, url string) (body []byte, statusCode int, err error)
}
