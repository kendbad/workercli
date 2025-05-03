package proxy

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type Checker struct {
	logger *utils.Logger
}

func NewChecker(logger *utils.Logger) *Checker {
	return &Checker{logger: logger}
}

func (c *Checker) CheckProxy(proxy model.Proxy, checkURL string) (model.ProxyResult, error) {
	result := model.ProxyResult{Proxy: proxy}

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
		RetryIf: func(req *fasthttp.Request) bool {
			return false // Disable retries for now; can enable for specific errors
		},
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(checkURL)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.DoTimeout(req, resp, 10*time.Second)
	if err != nil {
		result.Status = "Failed"
		result.Error = fmt.Sprintf("Không thể gửi request: %v", err)
		c.logger.Errorf("Proxy %s://%s:%s failed: %v", proxy.Protocol, proxy.IP, proxy.Port, err)
		return result, err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		result.Status = "Failed"
		result.Error = fmt.Sprintf("Mã trạng thái: %d", resp.StatusCode())
		c.logger.Errorf("Proxy %s://%s:%s returned status: %d", proxy.Protocol, proxy.IP, proxy.Port, resp.StatusCode())
		return result, fmt.Errorf("mã trạng thái: %d", resp.StatusCode())
	}

	var ipResp struct {
		IP string `json:"ip"`
	}
	if err := json.Unmarshal(resp.Body(), &ipResp); err != nil {
		result.Status = "Failed"
		result.Error = fmt.Sprintf("Không thể giải mã JSON: %v", err)
		c.logger.Errorf("Proxy %s://%s:%s failed JSON decode: %v", proxy.Protocol, proxy.IP, proxy.Port, err)
		return result, err
	}

	result.Status = "Success"
	result.IP = ipResp.IP
	c.logger.Infof("Proxy %s://%s:%s trả về IP: %s", proxy.Protocol, proxy.IP, proxy.Port, result.IP)
	return result, nil
}
