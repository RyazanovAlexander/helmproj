package io

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GetRuntimeFilePath returns the full path to the file relative to the executable program.
func GetRuntimeFilePath(filePath string) (string, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	return strings.Replace(absPath, "\\ \\", "\\", 1), nil
}

// RemoveAll removes all files and folders along the given path.
func RemoveAll(path string) error {
	if len(path) == 0 {
		return errors.New("Invalid parameter value passed: param - 'path', len(path) = 0")
	}

	toPath, err := GetRuntimeFilePath(path)
	if err != nil {
		return err
	}

	return os.RemoveAll(toPath)
}

// WriteToFile writes text to a file with overwriting
func WriteToFile(filePath string, data string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		dir := filepath.Dir(filePath)
		os.MkdirAll(dir, 0700)
	}

	if err := ioutil.WriteFile(filePath, []byte(data), 0644); err != nil {
		return err
	}

	return nil
}
