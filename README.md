# WorkerCLI

WorkerCLI là một ứng dụng CLI được viết bằng Go, tập trung vào xử lý số lượng cực lớn task đồng thời với hiệu suất cao thông qua hệ thống worker đa luồng. Ứng dụng sử dụng Clean Architecture, tích hợp logger để ghi lại hoạt động và hỗ trợ giao diện TUI. WorkerCLI được thiết kế để tối ưu xử lý hàng loạt task nhanh chóng, dễ dàng mở rộng cho các tính năng như gửi HTTP request hoặc kiểm tra email.

## Tính năng chính

- **Worker đa luồng**: Hệ thống worker pool mạnh mẽ, hỗ trợ xử lý đồng thời số lượng lớn task với queue tối ưu.
- **Hiệu suất cao**: Tối ưu cho tốc độ và quy mô, phù hợp với khối lượng task lớn.
- **Logger tích hợp**: Ghi log hoạt động với các cấp độ (debug, info, error) và định dạng (text, JSON).
- **TUI tích hợp**: Hỗ trợ giao diện người dùng văn bản (dùng thư viện `tview` `bubbletea`) để hiển thị tiến độ và kết quả khi cần.
- **Cấu hình linh hoạt**: Sử dụng file YAML để cấu hình worker, input/output, và logger.
- **Clean Architecture**: Mã nguồn được thiết kế mục đích học tập, cho người mới tiếp cận Clean Architecture. Thiết kế mô-đun, dễ bảo trì và mở rộng.

## Yêu cầu

- Go 1.21 trở lên.
- Thư viện (tự động tải qua `go mod tidy`):
  - `github.com/sirupsen/logrus@master`
  - `gopkg.in/yaml.v3@master`
  - `github.com/valyala/fasthttp@master`
  - `github.com/charmbracelet/bubbles@master`
  - `github.com/charmbracelet/bubbletea@master`
  - `github.com/rivo/tview@master`
  - `github.com/fatih/color@master`

## Cài đặt

1. Clone repository:

   ```bash
   git clone https://github.com/<your-repo>/workercli.git
   cd workercli
   ```

2. Cài đặt phụ thuộc:

   ```bash
   go mod tidy
   ```

3. (Tùy chọn) Cài đặt toàn cục:

   ```bash
   go install github.com/<your-repo>/workercli@latest
   ```

## Cách sử dụng

1. **Chuẩn bị file đầu vào**:

   - Tạo file `input/tasks.txt` với danh sách task (mỗi dòng là một task):

     ```
     task-data-1
     task-data-2
     task-data-3
     ```

2. **Cấu hình**:

   - Chỉnh sửa các file trong thư mục `configs/`:
     - `input.yaml`: Đường dẫn file đầu vào.
     - `output.yaml`: Đường dẫn file đầu ra.
     - `worker.yaml`: Số lượng worker và kích thước queue (tùy chỉnh cho khối lượng lớn).
     - `logger.yaml`: Cấp độ log, định dạng, và đầu ra (console hoặc file).

   Ví dụ `configs/worker.yaml`:

   ```yaml
   workers: 100
   queue_size: 10000
   ```

   Ví dụ `configs/logger.yaml`:

   ```yaml
   level: "info"
   output: "console"
   file_path: "logs/app.log"
   format: "text"
   ```

3. **Chạy ứng dụng**:

   - Chạy ở chế độ CLI (mặc định):

     ```bash
     go run cmd/workercli/main.go -task
     go run cmd/workercli/main.go -proxy
     ```

   - Chạy với TUI:

     ```bash
     go run cmd/workercli/main.go -proxy -tui tview
     go run cmd/workercli/main.go -task -tui tview
     go run cmd/workercli/main.go -proxy -tui bubbletea
     go run cmd/workercli/main.go -task -tui bubbletea
     ```

   Ứng dụng sẽ:

   - Đọc task từ file đầu vào.
   - Phân phối task qua worker pool với hiệu suất cao.
   - Ghi kết quả vào file đầu ra (nếu cấu hình).
   - Hiển thị tiến độ/kết quả trong TUI (nếu bật).
   - Ghi log hoạt động vào console hoặc file.

## Cấu trúc thư mục Clean Architecture

```bash
workercli/
├── cmd/                          # Điểm vào của ứng dụng
│   └── workercli/
│       └── main.go               # Điểm vào CLI
├── internal/                     # Logic cốt lõi (Clean Architecture)
│   ├── config/                   # Quản lý cấu hình
│   ├── domain/                   # Mô hình và logic nghiệp vụ
│   ├── usecase/                  # Các trường hợp sử dụng
│   ├── interface/                # Bộ điều hợp (input, worker)
│   └── infrastructure/           # Triển khai dịch vụ ngoài
├── configs/                      # File cấu hình YAML
├── pkg/                          # Công cụ dùng chung (logger, tui)
├── input/                        # File đầu vào
├── output/                       # File đầu ra
├── logs/                         # File log
├── go.mod                        # Module Go
├── go.sum                        # Checksum phụ thuộc
├── README.md                     # Tài liệu dự án
└── .gitignore                    # File bỏ qua git
```

Dưới đây là cấu trúc thư mục mới sau khi tổ chức lại TUI:
```bash
workercli/
├── cmd/
│   └── workercli/
│       └── main.go
├── configs/
│   ├── input.yaml
│   ├── output.yaml
│   ├── worker.yaml
│   ├── logger.yaml
│   └── proxy.yaml
├── input/
│   ├── tasks.txt
│   └── proxy.txt
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── types.go
│   ├── domain/
│   │   ├── model/
│   │   │   ├── task.go
│   │   │   ├── result.go
│   │   │   └── proxy.go
│   │   │   └── config.go
│   │   ├── service/
│   │   │   ├── task.go
│   │   │   └── proxy.go
├── infrastructure/
│   └── tui/
│       ├── types.go          # Cấu hình chung (RendererConfig, ComponentStyle)
│       ├── bubbletea/
│       │   ├── models/
│       │   │   └── status.go  # TUIRow, TUIViewState cho bubbletea
│       │   ├── renderer.go
│       │   ├── proxy_renderer.go
│       │   └── components/
│       │       ├── table.go
│       │       └── status.go
│       ├── tview/
│       │   ├── models/
│       │   │   └── model.go  # TUIRow, TUIViewState cho tview
│       │   ├── renderer.go
│       │   ├── proxy_renderer.go
│       │   └── components/
│       │       ├── layout.go
│       │       └── form.go
│       ├── termui/           # Thư viện TUI mới
│       │   ├── models/
│       │   │   └── model.go
│       │   ├── renderer.go
│       │   └── components/
│       │       ├── table.go
│       │       └── chart.go
│       ├── coordinator.go    # Điều phối renderer
│       └── factory.go        # Triển khai RendererFactory
│   ├── interface/
│   │   ├── input/
│   │   │   ├── file_reader.go
│   │   │   └── parser.go
│   │   ├── worker/
│   │   │   ├── worker.go
│   │   │   └── pool.go
│   │   ├── proxy/
│   │   │   ├── reader.go
│   │   │   └── checker.go
│   │   ├── tui/
│   │   │   ├── factory.go
│   │   │   ├── renderer.go         # Giao diện trừu tượng cho renderer
│   │   │   └── types.go           # Các kiểu dữ liệu cho TUI
│   ├── usecase/
│   │   ├── batch_task.go
│   │   └── proxy_check.go    
├── logs/
│   └── app.log
├── output/
│   ├── results.txt
│   └── proxy_results.txt
├── pkg/
│   ├── utils/
│   │   ├── logger.go
│   │   └── utils.go
├── go.mod
├── go.sum
├── README.md
└── .gitignore
```

## Mở rộng

WorkerCLI được thiết kế để dễ dàng thêm các tính năng mới:

- **Gửi HTTP request**: Thêm `RequestSender` vào `domain/service/` và triển khai trong `infrastructure/`.
- **Kiểm tra email**: Tạo `EmailChecker` và các mô hình liên quan (như `EmailAccount`).

## Ví dụ Log đầu ra

```
INFO[2025-05-02T12:00:00+07:00] Ứng dụng WorkerCLI khởi động
INFO[2025-05-02T12:00:01+07:00] Đọc được 1000 task từ file
INFO[2025-05-02T12:00:01+07:00] Khởi động pool với 100 worker
DEBUG[2025-05-02T12:00:02+07:00] Worker 1 nhận task task-1
DEBUG[2025-05-02T12:00:02+07:00] Xử lý task task-1 với dữ liệu: task-data-1
INFO[2025-05-02T12:00:02+07:00] Kết quả task task-1: success
DEBUG[2025-05-02T12:00:02+07:00] Worker 1 hoàn thành task task-1 với trạng thái success
INFO[2025-05-02T12:00:02+07:00] Hoàn thành xử lý 1000 task
```

## Góp ý

Nếu bạn có ý tưởng hoặc gặp vấn đề, vui lòng mở issue trên GitHub repository.

## Giấy phép

MIT License
