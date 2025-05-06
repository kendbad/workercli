package httpclient

import (
	"workercli/internal/domain/model"
	httpclient "workercli/internal/infrastructure/httpclient"
	"workercli/pkg/utils"
)

// HTTPClient defines the interface for sending HTTP requests.
type HTTPClient interface {
	DoRequest(proxy model.Proxy, duongDan string) ([]byte, int, error)
}

// Factory creates an HTTPClient based on the client type.
func NewHTTPClient(loaiKetNoi string, boGhiNhatKy *utils.Logger) HTTPClient {
	switch loaiKetNoi {
	case "fasthttp":
		return httpclient.NewFastHTTPClient(boGhiNhatKy)
	case "nethttp":
		return httpclient.NewNetHTTPClient(boGhiNhatKy)
	default:
		boGhiNhatKy.Warnf("Loại client không xác định %s, sử dụng fasthttp", loaiKetNoi)
		return httpclient.NewFastHTTPClient(boGhiNhatKy)
	}
}
