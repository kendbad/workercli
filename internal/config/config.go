package config

import (
	"os"
	"path/filepath"
	"workercli/pkg/utils"

	"gopkg.in/yaml.v3"
)

type ProxyConfig struct {
	FilePath string `yaml:"file_path"`
	CheckURL string `yaml:"check_url"`
}

type Config struct {
	Worker WorkerConfig `yaml:"worker"`
	Input  InputConfig  `yaml:"input"`
	Output OutputConfig `yaml:"output"`
	Proxy  ProxyConfig  `yaml:"proxy"`
}

type WorkerConfig struct {
	Workers   int `yaml:"workers"`
	QueueSize int `yaml:"queue_size"`
}

type InputConfig struct {
	FilePath string `yaml:"file_path"`
}

type OutputConfig struct {
	FilePath string `yaml:"file_path"`
}

func Load(configDir string) (*Config, error) {
	configDir = utils.AutoPath(configDir)
	cfg := &Config{}
	files := []string{"worker.yaml", "input.yaml", "output.yaml", "proxy.yaml"}

	for _, file := range files {
		data, err := os.ReadFile(filepath.Join(configDir, file))
		if err != nil {
			return nil, err
		}

		switch file {
		case "worker.yaml":
			err = yaml.Unmarshal(data, &cfg.Worker)
		case "input.yaml":
			err = yaml.Unmarshal(data, &cfg.Input)
		case "output.yaml":
			err = yaml.Unmarshal(data, &cfg.Output)
		case "proxy.yaml":
			err = yaml.Unmarshal(data, &cfg.Proxy)
		}
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
