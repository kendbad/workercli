package service

import (
	"workercli/internal/domain/model"
)

type BoKiemTraProxy interface {
	KiemTraProxy(proxy model.Proxy, duongDanKiemTra string) (ip string, trangThai string, err error)
}

// Interface này mô tả rằng "chúng tôi cần một thứ có thể kiểm tra proxy, bằng cách đưa vào proxy và checkURL, và nhận về địa chỉ IP và trạng thái".
// Nhưng không quan tâm ai và như thế nào sẽ thực thi.
// Điều này giúp cho việc thay đổi phương thức kiểm tra proxy trong tương lai sẽ dễ dàng hơn.
// Ví dụ: chúng ta có thể thay đổi phương thức kiểm tra proxy bằng cách sử dụng một thư viện khác, hoặc sử dụng một API khác.
// Chúng ta không cần phải thay đổi code của chính mình.
// Chúng ta chỉ cần tạo một struct mới và implement interface BoKiemTraProxy.
// Sau đó, chúng ta có thể sử dụng struct mới đó để kiểm tra proxy.
