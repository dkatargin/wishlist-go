package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type AppConfigStruct struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Worker struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	}
	Telegram struct {
		BotToken string `yaml:"bot_token"`
	} `yaml:"telegram"`
	RabbitMQ struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Vhost    string `yaml:"vhost"`
	} `yaml:"rabbitmq"`
	Sentry struct {
		DSN         string `yaml:"dsn"`
		Environment string `yaml:"environment"`
		Release     string `yaml:"release"`
	}
	Logging struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	}
}

var (
	once    sync.Once
	loadErr error
	Config  *AppConfigStruct
)

func LoadConfigFile(configPath string) error {
	once.Do(func() {

		data, err := os.ReadFile(configPath)
		if err != nil {
			loadErr = fmt.Errorf("failed to read config file: %w", err)
			return
		}

		if err := yaml.Unmarshal(data, &Config); err != nil {
			loadErr = fmt.Errorf("failed to unmarshal config: %w", err)
			return
		}
	})

	return loadErr
}
