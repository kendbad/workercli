package proxy

import (
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// BoDocProxy định nghĩa giao diện để đọc danh sách proxy.
// Theo Clean Architecture, interface này là một port trong tầng adapter
// giúp định nghĩa cách tương tác với nguồn dữ liệu bên ngoài
type Reader interface {
	ReadProxies(nguon string) ([]model.Proxy, error)
}

// BoDocProxy là bộ điều hợp (adapter) để đọc danh sách proxy.
// Tuân thủ nguyên tắc Dependency Inversion, lớp này phụ thuộc vào interface
// thay vì các triển khai cụ thể
type BoDocProxy struct {
	boGhiNhatKy *utils.Logger // boGhiNhatKy: bộ ghi nhật ký
	boDoc       Reader        // boDoc: bộ đọc - triển khai cụ thể (ví dụ: bộ đọc tệp tin)
	BoDocMock   Reader        // Mock object dùng cho việc testing
}

// TaoBoDocProxy tạo một bộ điều hợp (adapter) mới để đọc danh sách proxy.
// Áp dụng Dependency Injection để tiêm phụ thuộc
func NewProxyReader(boGhiNhatKy *utils.Logger, boDoc Reader) *BoDocProxy {
	return &BoDocProxy{
		boGhiNhatKy: boGhiNhatKy,
		boDoc:       boDoc,
	}
}

// ReadProxies đọc danh sách proxy từ một nguồn cụ thể.
// Phương thức này ủy quyền việc đọc cho triển khai cụ thể (boDoc)
func (r *BoDocProxy) ReadProxies(nguon string) ([]model.Proxy, error) {
	// Nếu có mock object, sử dụng nó trong test
	if r.BoDocMock != nil {
		return r.BoDocMock.ReadProxies(nguon)
	}

	danhSachProxy, err := r.boDoc.ReadProxies(nguon)
	if err != nil {
		r.boGhiNhatKy.Errorf("Không thể đọc danh sách proxy từ %s: %v", nguon, err)
		return nil, err
	}
	r.boGhiNhatKy.Infof("Đã đọc %d proxy từ %s", len(danhSachProxy), nguon)
	return danhSachProxy, nil
}
