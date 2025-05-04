package proxy

import (
	"fmt"
	"strings"
	"workercli/internal/domain/model"
)

// Proxy represents a proxy configuration
type Proxy struct {
	Protocol string
	IP       string
	Port     string
}

// ParseProxy parses a proxy string (e.g., "http://1.2.3.4:8080") into a Proxy struct
func ParseProxy(proxyStr string) (model.Proxy, error) {
	parts := strings.SplitN(proxyStr, "://", 2)
	if len(parts) != 2 {
		return model.Proxy{}, fmt.Errorf("invalid proxy format: %s", proxyStr)
	}

	addrParts := strings.Split(parts[1], ":")
	if len(addrParts) != 2 {
		return model.Proxy{}, fmt.Errorf("invalid proxy address: %s", parts[1])
	}
	return model.Proxy{Protocol: parts[0], IP: addrParts[0], Port: addrParts[1]}, nil
}
