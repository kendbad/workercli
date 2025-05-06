package service

import "workercli/internal/domain/model"

// BoXuLyTacVu định nghĩa giao diện để xử lý một tác vụ đơn lẻ.
// Theo nguyên tắc Dependency Inversion Principle trong Clean Architecture,
// domain layer định nghĩa các interface mà các layer bên ngoài cần triển khai
type BoXuLyTacVu interface {
	XuLyTacVu(tacVu model.TacVu) (model.KetQua, error)
}

// Interface này mô tả rằng "chúng tôi cần một thứ có thể xử lý tác vụ, bằng cách đưa vào tác vụ và nhận về KetQua".
// Nhưng không quan tâm ai và như thế nào sẽ thực thi.
// Điều này giúp cho việc thay đổi phương thức xử lý tác vụ trong tương lai sẽ dễ dàng hơn.
//
// Trong Clean Architecture, interface này là một "port" hoặc "cổng giao tiếp"
// giữa domain layer và các layer bên ngoài như usecase layer và infrastructure layer.
// Bằng cách định nghĩa interface trong domain layer, chúng ta đảm bảo rằng:
// 1. Domain logic không phụ thuộc vào triển khai cụ thể
// 2. Các layer bên ngoài phải tuân thủ quy tắc do domain layer đặt ra
// 3. Khả năng thay thế các triển khai mà không ảnh hưởng đến domain logic
//
// Ví dụ: chúng ta có thể thay đổi phương thức xử lý tác vụ bằng cách sử dụng một thư viện khác, hoặc sử dụng một API khác.
// Chúng ta không cần phải thay đổi code của chính mình.
// Chúng ta chỉ cần tạo một struct mới và implement interface BoXuLyTacVu.
// Sau đó, chúng ta có thể sử dụng struct mới đó để xử lý tác vụ.
