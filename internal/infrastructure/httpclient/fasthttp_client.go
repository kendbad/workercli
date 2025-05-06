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
	boGhiNhatKy *utils.Logger
}

func NewFastHTTPClient(boGhiNhatKy *utils.Logger) *FastHTTPClient {
	return &FastHTTPClient{boGhiNhatKy: boGhiNhatKy}
}

func (c *FastHTTPClient) DoRequest(trungGian model.TrungGian, duongDan string) ([]byte, int, error) {
	client := &fasthttp.Client{
		Dial: func(addr string) (conn net.Conn, err error) {
			switch trungGian.GiaoDien {
			case "http", "https":
				dialer := fasthttpproxy.FasthttpHTTPDialer(fmt.Sprintf("%s:%s", trungGian.DiaChi, trungGian.Cong))
				conn, err = dialer(addr)
			case "socks5":
				dialer := fasthttpproxy.FasthttpSocksDialer(fmt.Sprintf("%s:%s", trungGian.DiaChi, trungGian.Cong))
				conn, err = dialer(addr)
			default:
				err = fmt.Errorf("giao thức không được hỗ trợ: %s", trungGian.GiaoDien)
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
	req.SetRequestURI(duongDan)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.DoTimeout(req, resp, 10*time.Second)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body(), resp.StatusCode(), nil
}
