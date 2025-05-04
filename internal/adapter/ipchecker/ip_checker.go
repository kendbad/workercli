package ipchecker

import (
	"workercli/internal/domain/model"
	ipchecker "workercli/internal/infrastructure/ipchecker"
	"workercli/pkg/utils"
)

// IPChecker defines the interface for checking IP through a proxy.
type IPChecker interface {
	CheckIP(proxy model.Proxy, checkURL string) (string, int, error)
}

// IPCheckerAdapter is the adapter for IP checking.
type IPCheckerAdapter struct {
	logger  *utils.Logger
	checker IPChecker // Specific implementation (e.g., API checker)
}

func NewIPChecker(clientType string, logger *utils.Logger) *IPCheckerAdapter {
	checker := ipchecker.NewAPIChecker(clientType, logger)
	return &IPCheckerAdapter{
		logger:  logger,
		checker: checker,
	}
}

func (c *IPCheckerAdapter) CheckIP(proxy model.Proxy, checkURL string) (string, int, error) {
	return c.checker.CheckIP(proxy, checkURL)
}
