.PHONY: build clean run test test-v test-unit test-integration format lint release docker-build

# Biến môi trường
APP_NAME=workercli
VERSION=1.0.0
BUILD_DIR=output
MAIN_PKG=cmd/workercli/main.go

# Phiên bản Go
GO_VERSION=1.23.0

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)

build-all: clean
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_PKG)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_PKG)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_PKG)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(MAIN_PKG)

clean:
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)

run:
	go run $(MAIN_PKG)

run-proxy:
	go run $(MAIN_PKG) -proxy

run-task:
	go run $(MAIN_PKG) -task

run-tui:
	go run $(MAIN_PKG) -proxy -tui tview

test:
	go test ./...

test-v:
	go test -v ./...

test-unit:
	go test ./internal/...

test-integration:
	go test ./test/integration/...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

format:
	go fmt ./...

lint:
	go vet ./...
	@if command -v golint > /dev/null; then \
		golint ./...; \
	else \
		echo "golint not installed, skipping..."; \
	fi

release: build-all
	mkdir -p $(BUILD_DIR)/release
	cp README.md $(BUILD_DIR)/release/
	cp LICENSE $(BUILD_DIR)/release/
	cp -r configs/ $(BUILD_DIR)/release/configs/
	cd $(BUILD_DIR) && tar -czf release/$(APP_NAME)-$(VERSION).tar.gz $(APP_NAME)-*
	@echo "Release đã được tạo tại $(BUILD_DIR)/release"

docker-build:
	docker build -t $(APP_NAME):$(VERSION) .

# Cài đặt các công cụ phát triển
install-dev-tools:
	go install golang.org/x/lint/golint@latest 