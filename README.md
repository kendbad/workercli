# WorkerCLI

WorkerCLI là một ứng dụng CLI được viết bằng Go, tập trung vào xử lý số lượng cực lớn task đồng thời với hiệu suất cao thông qua hệ thống worker đa luồng. Ứng dụng sử dụng Clean Architecture, tích hợp logger để ghi lại hoạt động và hỗ trợ giao diện TUI. WorkerCLI được thiết kế để tối ưu xử lý hàng loạt task nhanh chóng, dễ dàng mở rộng cho các tính năng như gửi HTTP request hoặc kiểm tra email.

## Tính năng chính

- **Worker đa luồng**: Hệ thống worker pool mạnh mẽ, hỗ trợ xử lý đồng thời số lượng lớn task với queue tối ưu.
- **Hiệu suất cao**: Tối ưu cho tốc độ và quy mô, phù hợp với khối lượng task lớn.
- **Logger tích hợp**: Ghi log hoạt động với các cấp độ (debug, info, error) và định dạng (text, JSON).
- **TUI tích hợp**: Hỗ trợ giao diện người dùng văn bản (dùng thư viện `tview` `bubbletea`) để hiển thị tiến độ và kết quả khi cần.
- **Cấu hình linh hoạt**: Sử dụng file YAML để cấu hình worker, input/output, và logger.
- **Clean Architecture**: Mã nguồn được thiết kế mục đích học tập, cho người mới tiếp cận Clean Architecture. Thiết kế mô-đun, dễ bảo trì và mở rộng.

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

Dưới đây là cấu trúc thư mục chú thích chi tiết:
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

🧩 Tại sao có viewmodel.go?

Trong TUI, bạn không nên render trực tiếp từ domain model vì:

Domain thường chứa dữ liệu "thô".
UI cần hiển thị dữ liệu "thân thiện" hơn (icon, màu sắc, format text, phân trang).
👉 viewmodel.go là lớp chuyển đổi từ domain.Model → ViewModel để TUI dễ xử lý và hiển thị.

Trong Clean Architecture, tầng UI không nên xử lý dữ liệu thô trực tiếp từ domain.
viewmodel.go là cầu nối giúp:
Chuyển domain.Model → UI-friendly model (ViewModel)
Format data (icon, màu, trạng thái text)
Gom nhóm hoặc phân trang
Giúp tách rõ:
UseCase → ViewModel → Component (VD: Table, StatusBar)

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

## Cài đặt yêu cầu:

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

## Mở rộng

WorkerCLI được thiết kế để dễ dàng thêm các tính năng mới:

- **Gửi HTTP request**: Thêm `RequestSender` vào `domain/service/` và triển khai trong `infrastructure/`.
- **Kiểm tra email**: Tạo `EmailChecker` và các mô hình liên quan (như `EmailAccount`).

## Ví dụ Log đầu ra khi không sử dụng TUI.

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
