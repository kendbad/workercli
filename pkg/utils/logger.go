package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"

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
	// Lấy đường dẫn file thực thi
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	// Lấy thư mục chứa file thực thi (cmd/workercli)
	execDir := filepath.Dir(execPath)

	// Giả định project root là cha của cmd/workercli
	projectRoot := filepath.Dir(filepath.Dir(execDir))

	// Đường dẫn thư mục logs nằm cùng cấp với cmd/
	logsDir := filepath.Join(projectRoot, "logs")

	// Tạo thư mục nếu chưa tồn tại
	if err := os.MkdirAll(logsDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Tạo tên file log theo timestamp

	// timestamp := time.Now().Format("20060102_150405")
	timestamp := time.Now().Format("20060102") // chỉ lấy ngày

	logFileName := filepath.Join(logsDir, "log_"+timestamp+".log")

	log.Printf("Log file path: %s\n", logFileName)

	// Tạo file
	logFile, err := os.Create(logFileName)
	if err != nil {
		return nil, err
	}

	return logFile, nil
}
