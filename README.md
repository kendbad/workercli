# WorkerCLI

WorkerCLI là một ứng dụng CLI được viết bằng Go, tập trung vào xử lý số lượng cực lớn task đồng thời với hiệu suất cao thông qua hệ thống worker đa luồng. Ứng dụng sử dụng Clean Architecture, tích hợp logger để ghi lại hoạt động và hỗ trợ giao diện TUI. WorkerCLI được thiết kế để tối ưu xử lý hàng loạt task nhanh chóng, dễ dàng mở rộng cho các tính năng như gửi HTTP request hoặc kiểm tra email.

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/kendbad/workercli/ci.yml?branch=main)
![Docker Image Version](https://img.shields.io/docker/v/kendbad/workercli?label=docker)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/kendbad/workercli)
![Go version](https://img.shields.io/badge/Go-1.23-blue)
![License](https://img.shields.io/github/license/kendbad/workercli)

## Tính năng chính

- **Worker đa luồng**: Hệ thống worker pool mạnh mẽ, hỗ trợ xử lý đồng thời số lượng lớn task với queue tối ưu.
- **Hiệu suất cao**: Tối ưu cho tốc độ và quy mô, phù hợp với khối lượng task lớn.
- **Logger tích hợp**: Ghi log hoạt động với các cấp độ (debug, info, error) và định dạng (text, JSON).
- **TUI tích hợp**: Hỗ trợ giao diện người dùng văn bản (dùng thư viện `tview` `bubbletea`) để hiển thị tiến độ và kết quả khi cần.
- **Cấu hình linh hoạt**: Sử dụng file YAML để cấu hình worker, input/output, và logger.
- **Clean Architecture**: Mã nguồn được thiết kế mục đích học tập, cho người mới tiếp cận Clean Architecture. Thiết kế mô-đun, dễ bảo trì và mở rộng.
- **Error Handling**: Hệ thống xử lý lỗi toàn diện với các cơ chế retry cho task thất bại.
- **Docker Support**: Hỗ trợ đầy đủ Docker với multi-stage build để giảm kích thước image.
- **CI/CD Pipeline**: GitHub Actions tự động hóa quá trình build, test và phát hành.

## Cài đặt yêu cầu:

- Go 1.23 trở lên.
- Thư viện (tự động tải qua `go mod tidy`):
  - `github.com/sirupsen/logrus@master`
  - `gopkg.in/yaml.v3@master`
  - `github.com/valyala/fasthttp@master`
  - `github.com/charmbracelet/bubbles@master`
  - `github.com/charmbracelet/bubbletea@master`
  - `github.com/rivo/tview@master`
  - `github.com/fatih/color@master`

## Cài đặt

### Cài đặt từ source

1. Clone repository:

   ```bash
   git clone https://github.com/kendbad/workercli.git
   cd workercli
   ```

2. Cài đặt phụ thuộc:

   ```bash
   go mod tidy
   ```

3. Build ứng dụng:

   ```bash
   make build
   ```

4. (Tùy chọn) Cài đặt toàn cục:

   ```bash
   go install github.com/kendbad/workercli@latest
   ```

### Sử dụng Docker

1. Pull image từ Docker Hub:

   ```bash
   docker pull kendbad/workercli:latest
   ```

2. Hoặc build từ source:

   ```bash
   docker build -t workercli:latest .
   ```

3. Chạy với Docker:

   ```bash
   docker run -v $(pwd)/input:/var/lib/workercli/input \
              -v $(pwd)/output:/var/lib/workercli/output \
              -v $(pwd)/logs:/var/log/workercli \
              kendbad/workercli
   ```

### Tải bản release

Các bản release được đóng gói sẵn cho nhiều nền tảng khác nhau. Bạn có thể tải về từ trang [GitHub Releases](https://github.com/kendbad/workercli/releases).

## Cách sử dụng

1. **Chuẩn bị file đầu vào**:

   - Tạo file `input/tasks.txt` với danh sách task (mỗi dòng là một task):

     ```
     task-data-1
     task-data-2
     task-data-3
     ```

   - Hoặc `input/proxy.txt` với danh sách proxy (mỗi dòng một proxy):
   
     ```
     http://1.2.3.4:8080
     socks5://5.6.7.8:1080
     ```

2. **Cấu hình**:

   - Chỉnh sửa các file trong thư mục `configs/`:
     - `input.yaml`: Đường dẫn file đầu vào.
     - `output.yaml`: Đường dẫn file đầu ra.
     - `worker.yaml`: Số lượng worker và kích thước queue (tùy chỉnh cho khối lượng lớn).
     - `logger.yaml`: Cấp độ log, định dạng, và đầu ra (console hoặc file).
     - `proxy.yaml`: URL kiểm tra proxy và cấu hình timeout.

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
     # Kiểm tra proxy
     ./output/workercli -proxy
     
     # Xử lý task
     ./output/workercli -task
     
     # Chỉ định loại HTTP client (mặc định: nethttp)
     ./output/workercli -proxy -client fasthttp
     ```

   - Chạy với TUI:

     ```bash
     # Sử dụng giao diện tview
     ./output/workercli -proxy -tui tview
     
     # Sử dụng giao diện bubbletea
     ./output/workercli -task -tui bubbletea
     ```

   Ứng dụng sẽ:

   - Đọc task từ file đầu vào.
   - Phân phối task qua worker pool với hiệu suất cao.
   - Ghi kết quả vào file đầu ra (nếu cấu hình).
   - Hiển thị tiến độ/kết quả trong TUI (nếu bật).
   - Ghi log hoạt động vào console hoặc file.

## Xử lý lỗi và cơ chế Retry

WorkerCLI được trang bị cơ chế xử lý lỗi toàn diện:

- **Tự động retry**: Các task thất bại sẽ được thử lại theo cấu hình retry.
- **Timeout**: Mỗi task có thời gian timeout riêng để tránh blocking.
- **Graceful shutdown**: Xử lý tắt ứng dụng một cách an toàn, đảm bảo không mất dữ liệu.

Để cấu hình retry, chỉnh sửa file cấu hình tương ứng:

```yaml
# configs/retry.yaml (tạo mới nếu chưa có)
max_retries: 3
retry_delay: 2s
```

## Giám sát hiệu suất

WorkerCLI cung cấp các số liệu về hiệu suất thông qua cờ `-metrics`:

```bash
./output/workercli -proxy -metrics
```

Số liệu được ghi vào file log và có thể xuất ra stdout, bao gồm:
- Số lượng task đã xử lý
- Thời gian xử lý trung bình
- Tỷ lệ thành công/thất bại
- Độ trễ của queue

## Mở rộng

WorkerCLI được thiết kế để dễ dàng thêm các tính năng mới:

- **Gửi HTTP request**: Thêm `RequestSender` vào `domain/service/` và triển khai trong `infrastructure/`.
- **Kiểm tra email**: Tạo `EmailChecker` và các mô hình liên quan (như `EmailAccount`).
- **Tích hợp API mới**: Thêm adapter mới trong thư mục `adapter/`.

Xem thêm chi tiết về kiến trúc và hướng dẫn mở rộng trong [Readme.dev.md](Readme.dev.md).

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

## Đóng góp

Chúng tôi rất hoan nghênh mọi đóng góp! Vui lòng làm theo các bước sau:

1. Fork repository
2. Tạo branch mới (`git checkout -b feature/amazing-feature`)
3. Commit thay đổi của bạn (`git commit -m 'Add some amazing feature'`)
4. Push lên branch (`git push origin feature/amazing-feature`)
5. Mở Pull Request

Xem [CONTRIBUTING.md](CONTRIBUTING.md) để biết thêm chi tiết.

## Lưu ý phiên bản

Dự án này sử dụng [SemVer](http://semver.org/) để quản lý phiên bản.

## Giấy phép

Dự án này được phân phối dưới Giấy phép MIT. Xem [LICENSE](LICENSE) để biết thêm thông tin.
