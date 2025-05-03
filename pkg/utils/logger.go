package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Logger struct {
	*logrus.Logger
}

type Config struct {
	Level    string `yaml:"level"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"file_path"`
	Format   string `yaml:"format"`
}

func NewLogger(configPath string) (*Logger, error) {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()

	switch cfg.Level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	if cfg.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	if cfg.Output == "file" && cfg.FilePath != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
			return nil, err
		}
		file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(file)
	} else {
		logger.SetOutput(os.Stdout)
	}

	return &Logger{Logger: logger}, nil
}

func loadConfig(configPath string) (*Config, error) {
	configPath = AutoPath(configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// createLogFile tạo tệp log trong thư mục logs với tên tệp dựa trên tên thư mục hiện tại
func CreateLogFile() (*os.File, error) {
	// Lấy thư mục hiện tại
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Lỗi khi lấy thư mục hiện tại: %w", err)
	}

	// Lấy tên thư mục cuối cùng trong đường dẫn hiện tại
	dirName := filepath.Base(currentDir)

	// Đường dẫn đến thư mục logs
	logDir := "logs"

	// Tạo tên tệp log dựa trên tên thư mục hiện tại
	logFileName := fmt.Sprintf("%s_app.log", dirName)

	// Kiểm tra xem thư mục logs đã tồn tại chưa, nếu chưa thì tạo
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, os.ModePerm) // Tạo thư mục logs nếu chưa tồn tại
		if err != nil {
			return nil, fmt.Errorf("Không thể tạo thư mục logs: %w", err)
		}
	}

	// Tạo đường dẫn đầy đủ đến tệp log trong thư mục logs
	logFilePath := filepath.Join(logDir, logFileName)

	// Tạo tệp log
	logFile, err := os.Create(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("Lỗi khi tạo file log: %w", err)
	}

	return logFile, nil
}
