# Đóng góp cho WorkerCLI

Cảm ơn bạn đã quan tâm đến việc đóng góp cho WorkerCLI! Dưới đây là hướng dẫn để giúp bạn bắt đầu.

## Quy trình đóng góp

1. Fork repository
2. Clone repository của bạn `git clone https://github.com/YOUR_USERNAME/workercli.git`
3. Tạo branch mới `git checkout -b feature/amazing-feature`
4. Commit thay đổi của bạn `git commit -m 'Add some amazing feature'`
5. Push lên branch `git push origin feature/amazing-feature`
6. Mở Pull Request

## Hướng dẫn viết mã nguồn

- Tuân thủ các quy tắc định dạng của Go (sử dụng `gofmt`)
- Viết unit test cho mọi tính năng mới
- Cập nhật tài liệu khi cần thiết
- Đảm bảo rằng tất cả các test đều pass trước khi gửi PR

## Quy trình review code

1. Mỗi PR sẽ được xem xét bởi ít nhất một maintainer.
2. Sửa đổi có thể được yêu cầu trước khi merge.
3. Sau khi được chấp thuận, PR sẽ được merge vào branch chính.

## Viết commit message

- Sử dụng thì hiện tại (ví dụ: "Add feature" không phải "Added feature")
- Dòng đầu tiên nên ngắn gọn (tối đa 50 ký tự)
- Tham khảo số issue nếu có (ví dụ: "Fix #123")
- Sử dụng các prefix thông dụng:
  - `feat:` Thêm tính năng mới
  - `fix:` Sửa lỗi
  - `docs:` Thay đổi tài liệu
  - `test:` Thêm hoặc sửa test
  - `refactor:` Tái cấu trúc mã nguồn
  - `style:` Thay đổi định dạng (không ảnh hưởng đến code)
  - `perf:` Cải thiện hiệu suất

## Báo cáo lỗi

Nếu bạn tìm thấy lỗi, vui lòng mở issue với các thông tin sau:
- Mô tả ngắn gọn về lỗi
- Các bước để tái hiện lỗi
- Hành vi mong đợi và hành vi thực tế
- Phiên bản Go và hệ điều hành
- Log lỗi (nếu có)

## Góp ý tính năng mới

Nếu bạn có ý tưởng cho tính năng mới:
1. Trước tiên hãy kiểm tra xem tính năng đó đã được đề xuất chưa
2. Mở issue mới với tiêu đề "Feature Request: [Tên tính năng]"
3. Mô tả chi tiết tính năng và lợi ích của nó

## Cộng đồng

Chúng tôi cam kết duy trì một cộng đồng thân thiện và tôn trọng. Vui lòng tuân thủ [Quy tắc ứng xử](CODE_OF_CONDUCT.md) của chúng tôi. 