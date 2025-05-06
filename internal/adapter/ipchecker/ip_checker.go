package ipchecker

import (
	"workercli/internal/domain/model"
	ipchecker "workercli/internal/infrastructure/ipchecker"
	"workercli/pkg/utils"
)

// IPChecker defines the interface for checking IP through a proxy.
type IPChecker interface {
	CheckIP(proxy model.Proxy, duongDanKiemTra string) (string, int, error)
}

// IPCheckerAdapter is the adapter for IP checking.
type IPCheckerAdapter struct {
	boGhiNhatKy *utils.Logger
	boKiemTra   IPChecker // Specific implementation (e.g., API checker)
}

func NewIPChecker(loaiKetNoi string, boGhiNhatKy *utils.Logger) *IPCheckerAdapter {
	boKiemTra := ipchecker.NewAPIChecker(loaiKetNoi, boGhiNhatKy)
	return &IPCheckerAdapter{
		boGhiNhatKy: boGhiNhatKy,
		boKiemTra:   boKiemTra,
	}
}

func (c *IPCheckerAdapter) CheckIP(proxy model.Proxy, duongDanKiemTra string) (string, int, error) {
	return c.boKiemTra.CheckIP(proxy, duongDanKiemTra)
}
