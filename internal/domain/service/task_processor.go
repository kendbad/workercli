package service

import "workercli/internal/domain/model"

type BoXuLyTacVu interface {
	XuLyTacVu(tacVu model.TacVu) (model.KetQua, error)
}

// Interface này mô tả rằng "chúng tôi cần một thứ có thể xử lý tác vụ, bằng cách đưa vào tác vụ và nhận về KetQua".
// Nhưng không quan tâm ai và như thế nào sẽ thực thi.
// Điều này giúp cho việc thay đổi phương thức xử lý tác vụ trong tương lai sẽ dễ dàng hơn.
// Ví dụ: chúng ta có thể thay đổi phương thức xử lý tác vụ bằng cách sử dụng một thư viện khác, hoặc sử dụng một API khác.
// Chúng ta không cần phải thay đổi code của chính mình.
// Chúng ta chỉ cần tạo một struct mới và implement interface BoXuLyTacVu.
// Sau đó, chúng ta có thể sử dụng struct mới đó để xử lý tác vụ.
