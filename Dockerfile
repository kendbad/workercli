# Sử dụng Go 1.23 làm base image
FROM golang:1.23-alpine AS builder

# Cài đặt các công cụ cần thiết
RUN apk add --no-cache make git

# Thiết lập thư mục làm việc
WORKDIR /app

# Sao chép các file go.mod và go.sum đầu tiên để tận dụng cache
COPY go.mod go.sum ./

# Tải các dependency
RUN go mod download

# Sao chép toàn bộ source code
COPY . .

# Build ứng dụng
RUN make build

# Tạo image chạy cuối cùng (multi-stage build)
FROM alpine:latest

# Cài đặt các công cụ cần thiết cho production
RUN apk add --no-cache ca-certificates tzdata

# Tạo user không phải root để tăng tính bảo mật
RUN adduser -D -g '' appuser

# Sao chép binary từ stage builder
COPY --from=builder /app/output/workercli /usr/local/bin/

# Sao chép các file cấu hình
COPY --from=builder /app/configs/ /etc/workercli/configs/

# Tạo thư mục input và output, phân quyền cho appuser
RUN mkdir -p /var/lib/workercli/input /var/lib/workercli/output /var/log/workercli && \
    chown -R appuser:appuser /var/lib/workercli /var/log/workercli

# Chuyển sang user không phải root
USER appuser

# Thiết lập WORKDIR
WORKDIR /var/lib/workercli

# Volume cho data và logs
VOLUME ["/var/lib/workercli/input", "/var/lib/workercli/output", "/var/log/workercli"]

# Lệnh chạy mặc định
ENTRYPOINT ["workercli"]

# Tham số mặc định, có thể ghi đè khi chạy container
CMD ["-proxy"]

# Metadata
LABEL maintainer="Your Name <your.email@example.com>"
LABEL version="1.0.0"
LABEL description="WorkerCLI - Ứng dụng xử lý task đa luồng" 