package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// templatesDir конфиги строго в папке templates
const templatesDir = "templates/"

// LoadConfig загрузка конфига
func LoadConfig(path string) (*TemplateConfig, error) {
	fullPath := filepath.Join(templatesDir, path)
	bytes, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	cfg := new(TemplateConfig)

	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
