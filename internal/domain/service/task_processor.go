package service

import "workercli/internal/domain/model"

type TaskProcessor interface {
	ProcessTask(task model.Task) (model.Result, error)
}

// Interface này mô tả rằng "chúng tôi cần một thứ có thể xử lý task, bằng cách đưa vào task và nhận về Result".
// Nhưng không quan tâm ai và như thế nào sẽ thực thi.
// Điều này giúp cho việc thay đổi phương thức xử lý task trong tương lai sẽ dễ dàng hơn.
// Ví dụ: chúng ta có thể thay đổi phương thức xử lý task bằng cách sử dụng một thư viện khác, hoặc sử dụng một API khác.
// Chúng ta không cần phải thay đổi code của chính mình.
// Chúng ta chỉ cần tạo một struct mới và implement interface TaskProcessor.
// Sau đó, chúng ta có thể sử dụng struct mới đó để xử lý task.
