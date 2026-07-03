// Package approot определяет корневую папку с ресурсами относительно бинарника
package approot

import (
	"fmt"
	"os"
	"path/filepath"
)

var anchors = []string{"materials", "Rules"}

func Setup() (string, error) {
	var starts []string

	if exe, err := os.Executable(); err == nil {
		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			exe = resolved
		}
		starts = append(starts, filepath.Dir(exe))
	}

	if wd, err := os.Getwd(); err == nil {
		starts = append(starts, wd)
	}

	for _, start := range starts {
		if root, ok := findRoot(start); ok {
			if err := os.Chdir(root); err != nil {
				return "", fmt.Errorf("не удалось перейти в корень проекта %q: %w", root, err)
			}
			return root, nil
		}
	}

	return "", fmt.Errorf("не найдена корневая папка с ресурсами (искали %v рядом с бинарником и в рабочей директории)", anchors)
}

func findRoot(dir string) (string, bool) {
	for {
		if hasAnchor(dir) {
			return dir, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func hasAnchor(dir string) bool {
	for _, a := range anchors {
		if _, err := os.Stat(filepath.Join(dir, a)); err == nil {
			return true
		}
	}
	return false
}
