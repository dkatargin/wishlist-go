package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfigFile(configPath string) (*AppConfigStruct, error) {
	var config *AppConfigStruct
	var loadErr error
	data, err := os.ReadFile(configPath)
	if err != nil {
		loadErr = fmt.Errorf("failed to read config file: %w", err)
		return nil, loadErr
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		loadErr = fmt.Errorf("failed to unmarshal config: %w", err)
		return nil, loadErr
	}
	return config, loadErr

}
