package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// templatesDir конфиги строго в папке templates
const templatesDir = "templates/"

type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

type RealFS struct{}

func (RealFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// LoadConfig загрузка файла конфига
func LoadConfig(fs FileReader, path string) (*TemplateConfig, error) {
	fullPath := filepath.Join(templatesDir, path)
	bytes, err := fs.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	return ParseConfig(bytes)
}

// ParseConfig обработка конфига
func ParseConfig(bytes []byte) (*TemplateConfig, error) {
	cfg := new(TemplateConfig)

	err := json.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
