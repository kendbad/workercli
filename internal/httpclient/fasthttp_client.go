package httpclient

import (
	"fmt"
	"net"
	"time"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type FastHTTPClient struct {
	logger *utils.Logger
}

func NewFastHTTPClient(logger *utils.Logger) *FastHTTPClient {
	return &FastHTTPClient{logger: logger}
}

func (c *FastHTTPClient) DoRequest(proxy model.Proxy, url string) ([]byte, int, error) {
	client := &fasthttp.Client{
		Dial: func(addr string) (conn net.Conn, err error) {
			switch proxy.Protocol {
			case "http", "https":
				dialer := fasthttpproxy.FasthttpHTTPDialer(fmt.Sprintf("%s:%s", proxy.IP, proxy.Port))
				conn, err = dialer(addr)
			case "socks5":
				dialer := fasthttpproxy.FasthttpSocksDialer(fmt.Sprintf("%s:%s", proxy.IP, proxy.Port))
				conn, err = dialer(addr)
			default:
				err = fmt.Errorf("giao thức không được hỗ trợ: %s", proxy.Protocol)
			}
			return conn, err
		},
		ReadTimeout:         5 * time.Second,
		WriteTimeout:        5 * time.Second,
		MaxConnDuration:     10 * time.Second,
		MaxIdleConnDuration: 30 * time.Second,
		MaxConnsPerHost:     2,
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.DoTimeout(req, resp, 10*time.Second)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body(), resp.StatusCode(), nil
}
