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
	logger *utils.Logger
}

func NewNetHTTPClient(logger *utils.Logger) *NetHTTPClient {
	return &NetHTTPClient{logger: logger}
}

func (c *NetHTTPClient) DoRequest(proxy model.Proxy, urlStr string) ([]byte, int, error) {
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
	proxyAddr := fmt.Sprintf("%s:%s", proxy.IP, proxy.Port)
	switch proxy.Protocol {
	case "http", "https":
		proxyURL, err := url.Parse(fmt.Sprintf("%s://%s", proxy.Protocol, proxyAddr))
		if err != nil {
			c.logger.Errorf("Invalid proxy URL %s: %v", proxyAddr, err)
			return nil, 0, fmt.Errorf("invalid proxy URL: %v", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	case "socks5":
		dialer, err := goproxy.SOCKS5("tcp", proxyAddr, nil, &net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		})
		if err != nil {
			c.logger.Errorf("Failed to create SOCKS5 dialer for %s: %v", proxyAddr, err)
			return nil, 0, fmt.Errorf("failed to create SOCKS5 dialer: %v", err)
		}
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	default:
		c.logger.Errorf("Unsupported proxy protocol: %s", proxy.Protocol)
		return nil, 0, fmt.Errorf("unsupported proxy protocol: %s", proxy.Protocol)
	}

	// Tạo HTTP client
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	// Tạo request
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		c.logger.Errorf("Failed to create HTTP request: %v", err)
		return nil, 0, err
	}

	// Gửi request
	resp, err := client.Do(req)
	if err != nil {
		c.logger.Errorf("HTTP request failed: %v", err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Đọc response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Errorf("Failed to read response body: %v", err)
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
