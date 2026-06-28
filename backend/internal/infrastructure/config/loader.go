package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Load читает конфигурацию из переменных окружения (12-factor).
// Отсутствие required-переменной (пароли, bot token) — ошибка на старте.
func Load() (*AppConfigStruct, error) {
	cfg := &AppConfigStruct{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config from env: %w", err)
	}
	return cfg, nil
}
