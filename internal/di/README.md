# Thư mục DI (Dependency Injection)

Thư mục này chứa các thành phần liên quan đến Dependency Injection (DI) trong ứng dụng WorkerCLI.

## Mục đích

Dependency Injection là một kỹ thuật trong lập trình, trong đó một đối tượng nhận các đối tượng khác (dependencies) mà nó phụ thuộc. Mục đích chính của DI là:

1. **Tách biệt các thành phần**: Giúp giảm sự phụ thuộc giữa các thành phần, làm cho code dễ bảo trì hơn.
2. **Dễ dàng thay thế thành phần**: Bằng cách sử dụng interface, ta có thể dễ dàng thay thế triển khai cụ thể.
3. **Dễ dàng kiểm thử**: Có thể mock các dependencies khi kiểm thử.

## Cấu trúc

- **container.go**: Định nghĩa Container DI, quản lý và cung cấp tất cả các dependencies cho ứng dụng.
- **interfaces.go**: Định nghĩa các interfaces sử dụng trong container.

## Container DI

Container DI trong ứng dụng này quản lý:

- Cấu hình ứng dụng
- Logger
- Các thành phần dành cho kiểm tra proxy
- Các thành phần dành cho xử lý tác vụ
- Giao diện người dùng (TUI)

## Sử dụng Container DI

Container DI được sử dụng trong `main.go` như sau:

```go
// Tạo container DI
container := di.NewContainer()
if err := container.KhoiTao(thuMucCauHinh); err != nil {
    return nil, err
}

// Sử dụng container để lấy các dependencies
boGhiNhatKy := container.LayBoGhiNhatKy()
boKiemTra, err := container.ThietLapKiemTraProxy(loaiKetNoi, kieuGiaoDien)
```

## Interface

Các interface được định nghĩa để hỗ trợ DI:

- `ProxyReaderInterface`: Đọc danh sách proxy từ nguồn dữ liệu
- `ProxyCheckerInterface`: Kiểm tra proxy
- `TaskProcessorInterface`: Xử lý tác vụ
- `InputReaderInterface`: Đọc dữ liệu đầu vào
- `TUIRendererInterface`: Hiển thị giao diện người dùng

Các interface này giúp container DI có thể hoạt động mà không cần biết chi tiết triển khai cụ thể, tuân theo nguyên tắc Dependency Inversion (nguyên tắc D trong SOLID). 