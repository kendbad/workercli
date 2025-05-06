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
│         Interface          │ ← adapter (TUI, file, nguoiXuLy)
└────────────┬───────────────┘
             │ depends on
┌────────────▼───────────────┐
│          Usecase           │ ← orchestrate logic (TacVu xử lý như nào)
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
├── input/              # Dữ liệu đầu vào (tacVu, proxy)
├── output/             # Kết quả sau khi xử lý
├── logs/               # Ghi log hệ thống
├── pkg/                # Thư viện dùng lại
├── adapter/            # Logic chính của app (Clean Architecture)
│   ├── config/         # Load cấu hình từ YAML
│   ├── domain/         # Các model, interface cốt lõi (không phụ thuộc bên ngoài)
│   ├── usecase/        # Tầng điều phối nghiệp vụ
│   ├── adapter/        # Kết nối giữa domain và bên ngoài (file, TUI, proxy, nguoiXuLy)
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
│   ├── worker.yaml                       # Cấu hình nguoiXuLy/nhomXuLy
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
│   │   │   ├── task.go                   # Struct đại diện cho nhiệm vụ (TacVu)
│   │   │   ├── proxy.go                  # Struct đại diện proxy (Proxy)
│   │   │   ├── result.go                 # Kết quả xử lý tacVu hoặc proxy (KetQua, KetQuaProxy)
│   │   │   └── config.go                 # Struct cấu hình nội bộ
│   │   └── service/                      # Interface của các logic xử lý domain
│   │       ├── task_service.go          # Interface xử lý tacVu (BoXuLyTacVu)
│   │       └── proxy_service.go         # Interface xử lý proxy (BoKiemTraProxy)
│
│   ├── usecase/                          # Application logic: điều phối hành vi dựa trên yêu cầu từ adapter
│   │   ├── batch_task.go                # Use case xử lý danh sách tacVu (XuLyLoDongTacVu)
│   │   └── proxy_check.go              # Use case kiểm tra proxy (KiemTraProxy)
│
│   ├── adapter/                          # Adapter layer: xử lý giao tiếp vào/ra hệ thống
│   │   ├── input/                        # Đọc file đầu vào (tacVu, proxy,...)
│   │   │   ├── file_reader.go            # Đọc file txt
│   │   │   └── parser.go                 # Parse nội dung file
│   │   ├── proxy/                        # Giao tiếp với logic kiểm tra proxy
│   │   │   ├── reader.go                 # Đọc danh sách proxy
│   │   │   └── checker.go                # Gửi request kiểm tra proxy
│   │   ├── worker/                       # Tạo nhomXuLy, xử lý đồng thời
│   │   │   ├── pool.go                   # Quản lý nhomXuLy (NhomXuLy)
│   │   │   └── worker.go                 # Một nguoiXuLy đơn lẻ (NguoiXuLy)
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
│   │       │       ├── table.go          # Bảng hiển thị tacVu/proxy
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
│   ├── tasks.txt                          # Danh sách tacVu
│   └── proxy.txt                          # Danh sách proxy
│
├── output/                                # Kết quả xuất ra
│   ├── results.txt                        # Kết quả tacVu
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

Sau đó thêm httpclient, ipchecker, giống cách TUI phân bổ 2 layer: 

```bash
workercli/
├── cmd/
│   └── workercli/ 
│       └── main.go                       // Hàm main, khởi tạo hệ thống
├── configs/
│   ├── input.yaml                       // Cấu hình dữ liệu đầu vào
│   ├── output.yaml                      // Cấu hình xuất dữ liệu
│   ├── worker.yaml                      // Cấu hình nguoiXuLy/nhomXuLy
│   ├── logger.yaml                      // Cấu hình logger
│   └── proxy.yaml                       // Cấu hình kiểm tra proxy
├── internal/
│   ├── config/
│   │   ├── loader.go                    // Đọc và parse file YAML
│   │   └── model.go                     // Struct ánh xạ cấu hình
│   ├── domain/
│   │   ├── model/
│   │   │   ├── task.go                  // Struct TacVu
│   │   │   ├── proxy.go                 // Struct Proxy và ParseProxy
│   │   │   ├── result.go                // Struct KetQua và KetQuaProxy
│   │   │   └── config.go                // Struct cấu hình nội bộ
│   │   └── service/
│   │       ├── task_service.go          // Interface xử lý tacVu (BoXuLyTacVu)
│   │       └── proxy_service.go         // Interface xử lý proxy (BoKiemTraProxy)
│   ├── usecase/
│   │   ├── batch_task.go                // Use case xử lý danh sách tacVu (XuLyLoDongTacVu)
│   │   └── proxy_check.go               // Use case kiểm tra proxy (KiemTraProxy)
│   ├── adapter/
│   │   ├── input/
│   │   │   ├── file_reader.go           // Đọc file txt
│   │   │   └── parser.go                // Parse nội dung file
│   │   ├── proxy/
│   │   │   ├── reader.go                // Interface và logic đọc proxy
│   │   │   └── checker.go               // Interface và logic kiểm tra proxy (BoKiemTra)
│   │   ├── httpclient/
│   │   │   └── http_client.go           // Interface HTTPClient
│   │   ├── ipchecker/
│   │   │   └── ip_checker.go            // Interface IPChecker
│   │   ├── worker/
│   │   │   ├── pool.go                  // Quản lý nhomXuLy (NhomXuLy)
│   │   │   └── worker.go                // Một nguoiXuLy đơn lẻ (NguoiXuLy)
│   │   └── tui/
│   │       ├── factory.go               // Tạo renderer TUI
│   │       ├── renderer.go              // Interface renderer
│   │       ├── types.go                 // Kiểu dữ liệu chung cho TUI
│   │       ├── coordinator.go           // Điều phối TUI
│   │       ├── tui_factory.go           // Factory chọn renderer
│   │       └── config.go                // Cấu hình TUI
│   ├── infrastructure/
│   │   ├── task/
│   │   │   └── processor.go             // Xử lý tacVu (BoXuLy)
│   │   ├── httpclient/
│   │   │   ├── fasthttp_client.go       // Triển khai fasthttp
│   │   │   └── nethttp_client.go        // Triển khai net/http
│   │   ├── proxy/
│   │   │   ├── file_reader.go           // Triển khai đọc proxy từ file
│   │   │   └── ip_checker.go            // Triển khai kiểm tra proxy qua ipchecker (BoKiemTraIP)
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
│   ├── tasks.txt                        // Danh sách tacVu
│   └── proxy.txt                        // Danh sách proxy
├── output/
│   ├── results.txt                      // Kết quả tacVu
│   └── proxy_results.txt                // Kết quả kiểm tra proxy
├── logs/
│   └── app.log                          // Log ứng dụng
├── go.mod                               // Module Go
├── go.sum                               // Checksum dependencies
├── README.md                            // Tài liệu dự án
└── .gitignore                           // File bỏ qua git
```

🧠 Lưu ý về kiến trúc:

1. Domain Layer (Tầng miền):
   - Chứa Model và Service.
   - Model: `TacVu`, `KetQua`, `Proxy`, `KetQuaProxy` - đại diện cho khái niệm trong hệ thống.
   - Service: `BoXuLyTacVu`, `BoKiemTraProxy` - định nghĩa giao diện xử lý.
   - Không phụ thuộc vào bất kỳ framework hay thư viện bên ngoài.

2. Usecase Layer (Tầng ứng dụng):
   - `XuLyLoDongTacVu`, `KiemTraProxy` - điều phối các tác vụ.
   - Chỉ phụ thuộc vào Domain Layer.
   - Thực hiện logic nghiệp vụ, không quan tâm đến chi tiết hiển thị UI hay lưu trữ.

3. Adapter Layer (Tầng tiếp hợp):
   - Định nghĩa giao diện trừu tượng với thế giới bên ngoài.
   - `Reader`, `BoKiemTra`, `NhomXuLy`, `NguoiXuLy` - là các giao diện.
   - Chỉ phụ thuộc vào Domain và Usecase.

4. Infrastructure Layer (Tầng hạ tầng):
   - Triển khai cụ thể các giao diện từ Adapter.
   - `FileReader`, `BoKiemTraIP`, `FastHTTPClient`, `NetHTTPClient` - là các triển khai cụ thể.
   - Có thể phụ thuộc vào thư viện bên ngoài.

Việc chuyển đổi tên các thành phần từ tiếng Anh sang tiếng Việt giúp:
1. Thống nhất quy ước đặt tên trong toàn bộ dự án
2. Dễ hiểu hơn cho người phát triển Việt Nam
3. Tuân thủ các nguyên tắc Clean Architecture và duy trì tính rõ ràng, phân tách giữa các tầng

## Hướng dẫn phát triển

### 1. Thêm usecase mới

Ví dụ thêm kiểm tra proxy SOCKS5:

1. Tạo model trong domain/model
2. Định nghĩa interface trong domain/service
3. Thêm usecase mới trong usecase/
4. Thêm adapter thích hợp
5. Triển khai cụ thể trong infrastructure/

### 2. Thêm UI mới (ngoài TView & BubbleTea)

1. Tạo thư mục mới trong infrastructure/tui/
2. Triển khai Renderer và các component cần thiết
3. Cập nhật factory để chọn renderer mới

### 3. Thêm HTTPClient mới

1. Triển khai interface từ adapter/httpclient trong infrastructure/httpclient/
2. Cập nhật factory để chọn client mới

## Quy tắc đặt tên

- Model domain: `TacVu`, `KetQua`, `Proxy`, `KetQuaProxy`
- Interface: `BoXuLyTacVu`, `BoKiemTraProxy`, `BoKiemTra`
- Triển khai cụ thể: `BoXuLy`, `BoKiemTraIP`
- Biến: sử dụng lowerCamelCase trong tiếng Việt (`maTacVu`, `boGhiNhatKy`)
- Hằng số: sử dụng SNAKE_CASE (`MAX_SO_LUONG_NGUOI_XU_LY`)
- Tên file: snake_case (`ip_checker.go`, `proxy_reader.go`)
- Tên package: một từ snake_case (`ipchecker`, `httpclient`)

## Testing

Các test được tổ chức theo cấu trúc thư mục tương ứng với mã nguồn:

```bash
workercli/
├── test/
│   ├── unit/
│   │   ├── domain/
│   │   ├── usecase/
│   │   ├── adapter/
│   │   └── infrastructure/
│   └── integration/
│       ├── adapter_infrastructure/
│       └── usecase_adapter/
```

Kiểm thử đơn vị nên sử dụng mock để kiểm tra nhiều trường hợp khác nhau.
