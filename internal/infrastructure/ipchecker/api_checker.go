package ipchecker

import (
	"encoding/json"
	"fmt"
	"strings"
	"workercli/internal/adapter/httpclient"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

type APIChecker struct {
	ketNoi      httpclient.HTTPClient
	boGhiNhatKy *utils.Logger
}

func NewAPIChecker(loaiKetNoi string, boGhiNhatKy *utils.Logger) *APIChecker {
	ketNoi := httpclient.NewHTTPClient(loaiKetNoi, boGhiNhatKy)
	return &APIChecker{ketNoi: ketNoi, boGhiNhatKy: boGhiNhatKy}
}

func (c *APIChecker) CheckIP(trungGian model.TrungGian, duongDanKiemTra string) (string, int, error) {
	// Đảm bảo duongDanKiemTra là URL hợp lệ
	if !strings.HasPrefix(duongDanKiemTra, "http://") && !strings.HasPrefix(duongDanKiemTra, "https://") {
		duongDanKiemTra = "http://" + duongDanKiemTra
	}

	noiDung, maKetQua, err := c.ketNoi.DoRequest(trungGian, duongDanKiemTra)
	if err != nil {
		c.boGhiNhatKy.Errorf("Kiểm tra IP thất bại cho proxy %s://%s:%s: %v", trungGian.GiaoDien, trungGian.DiaChi, trungGian.Cong, err)
		return "", maKetQua, err
	}

	var ketQuaIP struct {
		IP string `json:"ip"`
	}
	if err := json.Unmarshal(noiDung, &ketQuaIP); err != nil {
		c.boGhiNhatKy.Errorf("Không thể giải mã phản hồi JSON: %v", err)
		return "", maKetQua, fmt.Errorf("Không thể giải mã JSON: %v", err)
	}

	return ketQuaIP.IP, maKetQua, nil
}
