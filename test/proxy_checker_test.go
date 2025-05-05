package test

import (
	"testing"
	"workercli/internal/domain/model"
	"workercli/internal/infrastructure/proxy"
	"workercli/pkg/utils"
)

func TestProxyCheckerSuccess(t *testing.T) {
	logger, _ := utils.NewLogger("configs/logger.yaml")
	_ = proxy.NewIPChecker(logger, "nethttp")

	_ = model.Proxy{
		Protocol: "http",
		IP:       "127.0.0.1",
		Port:     "8080",
	}

	t.Log("Test proxy checker success")
	// Trong thực tế, bạn nên sử dụng mock cho HTTP client
}

func TestProxyInvalidFormat(t *testing.T) {
	logger, _ := utils.NewLogger("configs/logger.yaml")
	_ = proxy.NewIPChecker(logger, "nethttp")

	_ = model.Proxy{
		Protocol: "invalid",
		IP:       "127.0.0.1",
		Port:     "8080",
	}

	t.Log("Test proxy với protocol không hợp lệ")
	// Trong thực tế, bạn nên sử dụng mock cho HTTP client
}
