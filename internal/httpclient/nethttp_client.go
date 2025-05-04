package httpclient

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

type NetHTTPClient struct {
	logger *utils.Logger
}

func NewNetHTTPClient(logger *utils.Logger) *NetHTTPClient {
	return &NetHTTPClient{logger: logger}
}

func (c *NetHTTPClient) DoRequest(proxy model.Proxy, url string) ([]byte, int, error) {
	proxyURL := fmt.Sprintf("%s://%s:%s", proxy.Protocol, proxy.IP, proxy.Port)
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          2,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	if proxy.Protocol == "http" || proxy.Protocol == "https" {
		client.Transport.(*http.Transport).Proxy = http.ProxyURL(mustParseURL(proxyURL))
	} else if proxy.Protocol == "socks5" {
		return nil, 0, fmt.Errorf("SOCKS5 not supported with net/http")
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

func mustParseURL(urlStr string) *url.URL {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(fmt.Sprintf("Invalid proxy URL: %v", err))
	}
	return u
}
