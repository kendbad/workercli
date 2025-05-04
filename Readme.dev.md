# Tài liệu phát triển: workercli

## 🌐 Mục tiêu kiến trúc

Dự án `workercli` được thiết kế theo Clean Architecture để dễ dàng mở rộng, bảo trì và phân chia công việc theo từng tầng rõ ràng. Module mẫu `ip_checker` sẽ là chuẩn kiến trúc để các module mới như `check_email`, `check_account`, v.v... tham khảo và phát triển dựa theo.

🧱 Clean Architecture là gì?

Clean Architecture là một mẫu thiết kế giúp phân tách rõ ràng các tầng trong phần mềm để:

Dễ mở rộng: thay giao diện, thay storage, hoặc logic mà không ảnh hưởng phần còn lại.
Dễ test: logic chính tách khỏi giao diện hoặc I/O.
Dễ phân chia nhóm: nhóm A làm UI, nhóm B làm core logic.

⚙️ Tầng phân chia:

```
┌────────────────────────────┐
│         Interface          │ ← adapter (TUI, file, worker)
└────────────┬───────────────┘
             │ depends on
┌────────────▼───────────────┐
│          Usecase           │ ← orchestrate logic (Task xử lý như nào)
└────────────┬───────────────┘
             │ depends on
┌────────────▼───────────────┐
│           Domain           │ ← business model & service
└────────────┬───────────────┘
             │
┌────────────▼───────────────┐
│       External Layer       │ ← infrastructure (lib cụ thể: bubbletea)
└────────────────────────────┘
```

## Cấu trúc thư mục Clean Architecture

```bash
workercli/
├── cmd/                # Điểm vào chính của app
├── configs/            # File cấu hình YAML
├── input/              # Dữ liệu đầu vào (task, proxy)
├── output/             # Kết quả sau khi xử lý
├── logs/               # Ghi log hệ thống
├── pkg/                # Thư viện dùng lại
├── adapter/            # Logic chính của app (Clean Architecture)
│   ├── config/         # Load cấu hình từ YAML
│   ├── domain/         # Các model, interface cốt lõi (không phụ thuộc bên ngoài)
│   ├── usecase/        # Tầng điều phối nghiệp vụ
│   ├── adapter/        # Kết nối giữa domain và bên ngoài (file, TUI, proxy, worker)
│   └── infrastructure/ # Cài đặt cụ thể TUI (bubbletea, tview,...)
└── README.md           # Tài liệu hướng dẫn
```

Dưới đây là cấu trúc thư mục chú thích chi tiết khi đã thêm TUI:

```bash
workercli/
├── cmd/                                  # Điểm khởi chạy ứng dụng
│   └── workercli/                        # Module chính
│       └── main.go                       # Hàm main, khởi tạo toàn bộ hệ thống
│
├── configs/                              # Các tệp cấu hình YAML của hệ thống
│   ├── input.yaml                        # Cấu hình cho dữ liệu đầu vào
│   ├── output.yaml                       # Cấu hình xuất dữ liệu
│   ├── worker.yaml                       # Cấu hình worker/pool
│   ├── logger.yaml                       # Cấu hình logger
│   └── proxy.yaml                        # Cấu hình kiểm tra proxy
│
├── internal/                             # Logic nội bộ của ứng dụng (theo Clean Architecture)
│   ├── config/                           # Loader config và model cấu hình
│   │   ├── loader.go                     # Hàm đọc và parse file YAML
│   │   └── model.go                      # Struct ánh xạ cấu hình
│
│   ├── domain/                           # Business domain: định nghĩa logic cốt lõi và giao diện (interface)
│   │   ├── model/                        # Các struct đại diện cho dữ liệu trong domain
│   │   │   ├── task.go                   # Struct đại diện cho nhiệm vụ
│   │   │   ├── proxy.go                  # Struct đại diện proxy
│   │   │   ├── result.go                 # Kết quả xử lý task hoặc proxy
│   │   │   └── config.go                 # Struct cấu hình nội bộ
│   │   └── service/                      # Interface của các logic xử lý domain
│   │       ├── task_service.go          # Interface xử lý task
│   │       └── proxy_service.go         # Interface xử lý proxy
│
│   ├── usecase/                          # Application logic: điều phối hành vi dựa trên yêu cầu từ adapter
│   │   ├── batch_task.go                # Use case xử lý danh sách task
│   │   └── proxy_check.go              # Use case kiểm tra proxy
│
│   ├── adapter/                          # Adapter layer: xử lý giao tiếp vào/ra hệ thống
│   │   ├── input/                        # Đọc file đầu vào (task, proxy,...)
│   │   │   ├── file_reader.go            # Đọc file txt
│   │   │   └── parser.go                 # Parse nội dung file
│   │   ├── proxy/                        # Giao tiếp với logic kiểm tra proxy
│   │   │   ├── reader.go                 # Đọc danh sách proxy
│   │   │   └── checker.go                # Gửi request kiểm tra proxy
│   │   ├── worker/                       # Tạo worker pool, xử lý đồng thời
│   │   │   ├── pool.go                   # Quản lý worker pool
│   │   │   └── worker.go                 # Một worker đơn lẻ
│   │   └── tui/                          # Giao diện dòng lệnh TUI (terminal UI) — định nghĩa interface trừu tượng
│   │       ├── factory.go                # Tạo renderer TUI phù hợp
│   │       ├── renderer.go               # Interface renderer
│   │       └── types.go                  # Kiểu dữ liệu chung cho TUI
│
│   ├── infrastructure/                   # Cài đặt chi tiết, dùng thư viện ngoài (UI framework, logging,...)
│   │   └── tui/                          # Hiện thực giao diện terminal UI theo nhiều thư viện khác nhau
│   │       ├── bubbletea/                # Cài đặt TUI bằng thư viện Bubbletea
│   │       │   ├── renderer.go
│   │       │   ├── proxy_renderer.go
│   │       │   ├── viewmodel.go
│   │       │   └── components/
│   │       │       ├── table.go          # Bảng hiển thị task/proxy
│   │       │       └── status.go         # Thanh trạng thái (status bar)
│   │       ├── tview/                    # Cài đặt TUI bằng thư viện Tview
│   │       │   ├── renderer.go
│   │       │   ├── proxy_renderer.go
│   │       │   ├── viewmodel.go
│   │       │   └── components/
│   │       │       ├── layout.go         # Layout TUI
│   │       │       └── form.go           # Form nhập liệu
│   │       ├── termui/                   # Cài đặt TUI bằng thư viện TermUI
│   │       │   ├── renderer.go
│   │       │   ├── viewmodel.go
│   │       │   └── components/
│   │       │       ├── table.go
│   │       │       └── chart.go          # Biểu đồ thống kê (nếu có)
│   │       ├── coordinator.go            # Điều phối TUI đang dùng (bubbletea, tview,...)
│   │       ├── factory.go                # Factory chọn renderer phù hợp
│   │       └── config.go                 # Cấu hình giao diện TUI
│
├── pkg/                                   # Thư viện dùng lại được (shared utility)
│   ├── utils/
│   │   ├── logger.go                      # Cấu hình logger chung
│   │   └── stringutil.go                  # Hàm xử lý chuỗi tiện ích
│   └── logger/                            # Tách riêng package logger nếu dùng phức tạp hơn
│       └── logger.go
│
├── input/                                 # Dữ liệu đầu vào (cho testing hoặc thực tế)
│   ├── tasks.txt                          # Danh sách task
│   └── proxy.txt                          # Danh sách proxy
│
├── output/                                # Kết quả xuất ra
│   ├── results.txt                        # Kết quả task
│   └── proxy_results.txt                  # Kết quả kiểm tra proxy
│
├── logs/                                  # Log file ứng dụng
│   └── app.log
│
├── go.mod                                 # Module Go
├── go.sum                                 # Checksum cho dependencies
├── README.md                              # Tài liệu giới thiệu dự án
└── .gitignore                             # Bỏ qua file không cần track bởi git
```

Sau đó thêm nhiều httpclient, ipchecker, giống cách TUI phân bổ 2 layer: 

```bash
workercli/
├── cmd/
│   └── workercli/ 
│       └── main.go                       // Hàm main, khởi tạo hệ thống
├── configs/
│   ├── input.yaml                       // Cấu hình dữ liệu đầu vào
│   ├── output.yaml                      // Cấu hình xuất dữ liệu
│   ├── worker.yaml                      // Cấu hình worker/pool
│   ├── logger.yaml                      // Cấu hình logger
│   └── proxy.yaml                       // Cấu hình kiểm tra proxy
├── internal/
│   ├── config/
│   │   ├── loader.go                    // Đọc và parse file YAML
│   │   └── model.go                     // Struct ánh xạ cấu hình
│   ├── domain/
│   │   ├── model/
│   │   │   ├── task.go                  // Struct Task
│   │   │   ├── proxy.go                 // Struct Proxy và ParseProxy
│   │   │   ├── result.go                // Struct Result
│   │   │   └── config.go                // Struct cấu hình nội bộ
│   │   └── service/
│   │       ├── task_service.go          // Interface xử lý task
│   │       └── proxy_service.go         // Interface xử lý proxy
│   ├── usecase/
│   │   ├── batch_task.go                // Use case xử lý danh sách task
│   │   └── proxy_check.go               // Use case kiểm tra proxy
│   ├── adapter/
│   │   ├── input/
│   │   │   ├── file_reader.go           // Đọc file txt
│   │   │   └── parser.go                // Parse nội dung file
│   │   ├── proxy/
│   │   │   ├── reader.go                // Interface và logic đọc proxy
│   │   │   └── checker.go               // Interface và logic kiểm tra proxy
│   │   ├── httpclient/
│   │   │   └── http_client.go           // Interface HTTPClient
│   │   ├── ipchecker/
│   │   │   └── ip_checker.go            // Interface IPChecker
│   │   ├── worker/
│   │   │   ├── pool.go                  // Quản lý worker pool
│   │   │   └── worker.go                // Một worker đơn lẻ
│   │   └── tui/
│   │       ├── factory.go               // Tạo renderer TUI
│   │       ├── renderer.go              // Interface renderer
│   │       ├── types.go                 // Kiểu dữ liệu chung cho TUI
│   │       ├── coordinator.go           // Điều phối TUI
│   │       ├── tui_factory.go           // Factory chọn renderer
│   │       └── config.go                // Cấu hình TUI
│   ├── infrastructure/
│   │   ├── task/
│   │   │   └── processor.go             // Xử lý task
│   │   ├── httpclient/
│   │   │   ├── fasthttp_client.go       // Triển khai fasthttp
│   │   │   └── nethttp_client.go        // Triển khai net/http
│   │   ├── proxy/
│   │   │   ├── file_reader.go           // Triển khai đọc proxy từ file
│   │   │   └── ip_checker.go            // Triển khai kiểm tra proxy qua ipchecker
│   │   ├── ipchecker/
│   │   │   └── api_checker.go           // Triển khai kiểm tra IP qua API
│   │   └── tui/
│   │       ├── bubbletea/
│   │       │   ├── renderer.go          // Triển khai Bubbletea
│   │       │   ├── proxy_renderer.go    // Renderer cho proxy
│   │       │   ├── viewmodel.go         // View model
│   │       │   └── components/
│   │       │       ├── table.go         // Bảng hiển thị
│   │       │       └── status.go        // Thanh trạng thái
│   │       ├── tview/
│   │       │   ├── renderer.go          // Triển khai Tview
│   │       │   ├── proxy_renderer.go    // Renderer cho proxy
│   │       │   ├── viewmodel.go         // View model
│   │       │   └── components/
│   │       │       ├── layout.go        // Layout TUI
│   │       │       └── form.go          // Form nhập liệu
│   └── pkg/
│       ├── utils/
│       │   ├── logger.go                // Cấu hình logger
│       │   └── stringutil.go            // Xử lý chuỗi
│       └── logger/
│           └── logger.go                // Package logger
├── input/
│   ├── tasks.txt                        // Danh sách task
│   └── proxy.txt                        // Danh sách proxy
├── output/
│   ├── results.txt                      // Kết quả task
│   └── proxy_results.txt                // Kết quả kiểm tra proxy
├── logs/
│   └── app.log                          // Log ứng dụng
├── go.mod                               // Module Go
├── go.sum                               // Checksum dependencies
├── README.md                            // Tài liệu dự án
└── .gitignore                           // File bỏ qua git
```

🧠 Lưu ý về kiến trúc:

```
Layer	Vai trò	Biết được tầng nào?
domain	Entity + Interface business logic thuần túy	KHÔNG biết gì về usecase/infra
usecase	Logic điều phối các hành động	Chỉ biết domain và adapter
adapter	Nhận input (CLI/HTTP/file), gửi tới usecase	Biết usecase và infra
infra	TUI, file, logging, network,...	KHÔNG biết gì về usecase/domain
```

📌 Gợi ý: Team mới vào chỉ cần đọc các mục sau

```
README.md: Hướng dẫn tổng quan.
cmd/workercli/main.go: Entry chính, từ đây hiểu flow tổng thể.
internal/usecase/: Hiểu các hành vi của ứng dụng.
internal/infrastructure/tui/: Biết đang dùng framework TUI nào.
adapter/: Biết dữ liệu vào/ra và worker pool xử lý thế nào.
```

✅ Gợi ý phân chia nhóm (5 người)

```bash
Thành viên	Vai trò chính	Phạm vi code chính
1. Usecase Master	Phát triển logic nghiệp vụ	internal/usecase/, internal/domain/model/, internal/domain/service/
2. Adapter Guru	Kết nối với nguồn dữ liệu & giao tiếp	internal/adapter/ (input, proxy, httpclient, ipchecker, worker)
3. Infra Hacker	Thực thi chi tiết kỹ thuật	internal/infrastructure/ (triển khai các interface: http, proxy, ip, tui)
4. UI/TUI Engineer	TUI hiển thị & phản hồi người dùng	internal/adapter/tui/, internal/infrastructure/tui/
5. Config & Cmd Builder	Cấu hình, bootstrap & glue code	cmd/, configs/, internal/config/, pkg/logger/, main.go
```

📌 Gợi ý mở rộng giống TUI (2-layer cho mỗi service)
Với mỗi service như: httpclient, ipchecker, chia theo 2 layer:

```bash
internal/
├── adapter/
│   └── httpclient/
│       └── http_client.go        // Interface và logics gọi từ usecase
├── infrastructure/
│   └── httpclient/
│       ├── fasthttp_client.go    // Implement cụ thể
│       └── nethttp_client.go     // Implement khác (hoặc mock test)
```

Áp dụng tương tự cho:

```bash
ipchecker → chia adapter/ipchecker và infrastructure/ipchecker
```

Nguyên tắc:
```bash
adapter/ dùng trong usecase và inject từ ngoài vào, còn infrastructure/ chứa các implement thực tế, có thể thay thế.
```

✅ Ưu điểm của cách chia này:
```bash
Tách biệt nhiệm vụ rõ ràng → dễ test, dễ debug, dễ onboarding.
Làm việc song song không đụng nhau → mỗi người chỉ cần giao tiếp qua interface.
Dễ mở rộng nhiều thư viện cùng lúc (giống BubbleTea, TermUI, Tview).
Đảm bảo Clean Architecture: usecase không biết gì về implement cụ thể.
```

## 🧠 Phân chia công việc cho 5 người

### 👤 1. Người 1 - HTTP Client

* Chịu trách nhiệm xây dựng các client HTTP chuẩn theo interface (`internal/adapter/httpclient/http_client.go`).
* Đã có: `fasthttp_client.go`, `nethttp_client.go`.
* Sẽ dùng lại client này cho các module như: `ipchecker`, `emailchecker`, v.v.
* Làm việc nhiều ở `internal/infrastructure/httpclient`.

### 👤 2. Người 2 - IP Checker

* Viết logic kiểm tra IP qua API có proxy (interface `IPChecker`).
* Tận dụng `HTTPClient` để gửi request, trả về thông tin IP nếu thành công.
* File chính: `internal/adapter/ipchecker/ip_checker.go` + `infrastructure/ipchecker/api_checker.go`.
* Đây là module mẫu để các checker sau tham khảo.

### 👤 3. Người 3 - Worker & Pool

* Xử lý logic worker pool, chia task, giới hạn goroutine theo config.
* File: `internal/adapter/worker/pool.go` và `worker.go`.
* Kết hợp với `usecase/proxy_check.go` để gán task tương ứng cho worker.

### 👤 4. Người 4 - UseCase & Domain

* Viết logic xử lý IP checker ở tầng `usecase`.
* Interface và DTO định nghĩa trong `domain/model`, `domain/service`.
* Cầu nối giữa worker và infrastructure.
* Đảm bảo không gọi thẳng HTTP ở tầng này, chỉ qua interface.

### 👤 5. Người 5 - Giao diện TUI

* Cấu hình, hiển thị kết quả proxy/ip checker ra màn hình.
* Làm việc ở `internal/adapter/tui`, có thể chọn `bubbletea` hoặc `tview`.
* Renderer sẽ gọi từ ViewModel → cập nhật real-time kết quả.

---

## 🔄 Dòng chảy dữ liệu (kiểm tra IP)

```
Main → Load Config → Init Worker Pool
     ↘ input.txt → Task → Worker ↘
        ↘ proxy.txt       ↘ ip_checker
             ↘ http_client ↘ API trả về IP
                       ↘ ghi kết quả
```

---

## 🔧 Kiến trúc mẫu: ip\_checker

### Interface: `internal/adapter/ipchecker/ip_checker.go`

```go
package ipchecker

import (
    "context"
    "workercli/internal/domain/model"
)

type IPChecker interface {
    Check(ctx context.Context, proxy model.Proxy) (*model.Result, error)
}
```

### Triển khai: `internal/infrastructure/ipchecker/api_checker.go`

```go
package ipchecker

import (
    "context"
    "workercli/internal/domain/model"
    "workercli/internal/adapter/httpclient"
)

type APIChecker struct {
    Client httpclient.HTTPClient
}

func NewAPIChecker(client httpclient.HTTPClient) *APIChecker {
    return &APIChecker{Client: client}
}

func (c *APIChecker) Check(ctx context.Context, proxy model.Proxy) (*model.Result, error) {
    req := model.HttpRequest{
        URL:    "https://api.ipify.org?format=json",
        Method: "GET",
        Proxy:  &proxy,
    }
    resp, err := c.Client.Do(ctx, req)
    if err != nil {
        return nil, err
    }
    return &model.Result{
        Proxy: proxy,
        Output: string(resp.Body),
    }, nil
}
```

---

## 🏗 Mô hình mở rộng (ví dụ email\_checker)

* Tạo interface `EmailChecker`
* Triển khai tương tự `APIChecker`
* Sử dụng lại `HTTPClient` và Worker Pool đã có
* Kết quả dùng lại `Result`, chỉ cần phân loại loại task đầu vào

---

## 🧪 Chạy thử kiểm tra IP

```bash
cd cmd/workercli
go run main.go -proxy -tui tview
```

---

## 🧼 Ghi chú Clean Code

* Tất cả code phải dùng interface khi gọi giữa các tầng
* Không gọi HTTP ở tầng usecase/domain
* Luôn log lỗi ra `logs/app.log`
* Mỗi module nên viết test riêng cho tầng infrastructure và usecase

---

> Mọi module mới sau này đều phải tham khảo `ip_checker` để giữ vững kiến trúc nhất quán 💡
