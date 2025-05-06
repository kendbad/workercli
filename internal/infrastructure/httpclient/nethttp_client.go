package httpclient

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"

	goproxy "golang.org/x/net/proxy"
)

type NetHTTPClient struct {
	boGhiNhatKy *utils.Logger
}

func NewNetHTTPClient(boGhiNhatKy *utils.Logger) *NetHTTPClient {
	return &NetHTTPClient{boGhiNhatKy: boGhiNhatKy}
}

func (c *NetHTTPClient) DoRequest(proxy model.Proxy, duongDan string) ([]byte, int, error) {
	// Tạo http.Transport với cấu hình chung
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          2,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// Xử lý proxy dựa trên giao thức
	diaChiProxy := fmt.Sprintf("%s:%s", proxy.DiaChi, proxy.Cong)
	switch proxy.GiaoDien {
	case "http", "https":
		proxyURL, err := url.Parse(fmt.Sprintf("%s://%s", proxy.GiaoDien, diaChiProxy))
		if err != nil {
			c.boGhiNhatKy.Errorf("URL proxy không hợp lệ %s: %v", diaChiProxy, err)
			return nil, 0, fmt.Errorf("URL proxy không hợp lệ: %v", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	case "socks5":
		dialer, err := goproxy.SOCKS5("tcp", diaChiProxy, nil, &net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		})
		if err != nil {
			c.boGhiNhatKy.Errorf("Không thể tạo kết nối SOCKS5 cho %s: %v", diaChiProxy, err)
			return nil, 0, fmt.Errorf("không thể tạo kết nối SOCKS5: %v", err)
		}
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	default:
		c.boGhiNhatKy.Errorf("Giao thức proxy không được hỗ trợ: %s", proxy.GiaoDien)
		return nil, 0, fmt.Errorf("giao thức proxy không được hỗ trợ: %s", proxy.GiaoDien)
	}

	// Tạo HTTP client
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	// Tạo request
	req, err := http.NewRequest(http.MethodGet, duongDan, nil)
	if err != nil {
		c.boGhiNhatKy.Errorf("Không thể tạo request HTTP: %v", err)
		return nil, 0, err
	}

	// Gửi request
	resp, err := client.Do(req)
	if err != nil {
		c.boGhiNhatKy.Errorf("Request HTTP thất bại: %v", err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Đọc response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.boGhiNhatKy.Errorf("Không thể đọc phản hồi: %v", err)
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
