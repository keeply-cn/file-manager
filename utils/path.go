package utils

import (
	"errors"
	"path/filepath"
	"strings"
)

var ErrAccessDenied = errors.New("access denied")

func SafePath(rootDir, userPath string) (string, error) {
	if userPath == "" {
		userPath = "/"
	}

	absRoot, _ := filepath.Abs(rootDir)
	absPath, _ := filepath.Abs(filepath.Join(rootDir, userPath))

	if !strings.HasPrefix(absPath, absRoot) {
		return "", ErrAccessDenied
	}

	return absPath, nil
}

func GetBaseName(path string) string {
	return filepath.Base(path)
}

func GetDir(path string) string {
	return filepath.Dir(path)
}
