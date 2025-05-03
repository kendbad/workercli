package utils

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
)

// TODO: Thêm các hàm tiện ích chung nếu cần

// AutoPath tự động xác định đường dẫn tới thư mục configs
func AutoPath(path string) string {
	// 1. Nếu có truyền qua cờ -config thì ưu tiên dùng
	var configDir string
	flag.StringVar(&configDir, path, "", "Config directory path")
	flag.Parse()
	if configDir != "" {
		return configDir + path
	}

	// 2. Nếu có ENV CONFIG_DIR thì dùng
	if env := os.Getenv("CONFIG_DIR"); env != "" {
		return env + path
	}

	// 3. Ngược lại, lấy đường dẫn mặc định: đi từ file hiện tại tới thư mục gốc
	_, currentFile, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(filepath.Dir(currentFile), "..", "..", path)
	rootDir = filepath.Clean(rootDir)
	return rootDir
}
