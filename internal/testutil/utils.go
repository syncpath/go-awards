// Package testutil утилиты для написания тестов
package testutil

import (
	"errors"
	"os"
	"path/filepath"
)

func FindPath(target string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		fullpath := filepath.Join(dir, target)
		if _, err := os.Stat(fullpath); err == nil {
			return fullpath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", errors.New("папка tables не найдена")
}
