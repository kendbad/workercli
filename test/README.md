# Cấu Trúc Test Theo Clean Architecture

Thư mục này chứa các test được tổ chức theo nguyên tắc Clean Architecture, phản ánh cấu trúc của codebase và cô lập các thành phần để kiểm tra hiệu quả.

## Cấu trúc thư mục

```
test/
├── unit/                   # Test đơn vị cho từng tầng
│   ├── domain/             # Test cho domain entities và business rules
│   │   ├── model/          # Test cho domain models 
│   │   └── service/        # Test cho domain services
│   ├── usecase/            # Test cho application use cases
│   ├── adapter/            # Test cho interface adapters
│   └── infrastructure/     # Test cho frameworks và drivers
└── integration/            # Test tích hợp giữa các tầng
    ├── adapter_infrastructure/  # Test tích hợp giữa adapter và infrastructure
    └── usecase_adapter/    # Test tích hợp giữa usecase và adapter
```

## Nguyên tắc test

### Unit Tests (Test đơn vị)

Unit test kiểm tra từng thành phần riêng lẻ với các dependencies được mock/stub. Mỗi tầng của Clean Architecture được test một cách độc lập:

1. **Domain Tests**: Kiểm tra các entities và business rules cốt lõi không phụ thuộc vào các framework hay thư viện bên ngoài.
   
2. **Usecase Tests**: Test các trường hợp sử dụng ứng dụng bằng cách mock các dependencies như repositories và services.
   
3. **Adapter Tests**: Kiểm tra các adapter chuyển đổi dữ liệu giữa use cases và thế giới bên ngoài.
   
4. **Infrastructure Tests**: Kiểm tra các implementation cụ thể của framework và giao tiếp với thế giới bên ngoài.

### Integration Tests (Test tích hợp)

Integration test kiểm tra việc tương tác giữa các thành phần thực tế trong các tầng khác nhau:

1. **Adapter-Infrastructure**: Kiểm tra tương tác giữa adapter và infrastructure, đảm bảo dữ liệu được chuyển đổi và xử lý chính xác.
   
2. **Usecase-Adapter**: Kiểm tra tương tác giữa use cases và các adapter, xác nhận luồng nghiệp vụ end-to-end.

## Quy ước đặt tên

- Tên file test nên kết thúc bằng `_test.go`
- Tên package test nên là `package_test` để tránh circular dependencies
- Tên function test nên bắt đầu bằng `Test` và mô tả chức năng đang được kiểm tra

## Mocking trong tests

Chúng tôi sử dụng các interface để dễ dàng mock các dependencies. Điều này phù hợp với nguyên tắc Dependency Inversion của Clean Architecture, cho phép các unit test hoạt động mà không cần các thành phần bên ngoài thực tế.

## Chạy tests

Để chạy tất cả các tests:

```
go test ./test/...
```

Để chạy unit tests:

```
go test ./test/unit/...
```

Để chạy integration tests:

```
go test ./test/integration/...
``` 