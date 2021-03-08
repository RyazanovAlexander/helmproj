/*
MIT License

Copyright The Helmproj Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package io

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"

	"github.com/RyazanovAlexander/helmproj/v1/config"
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

// WriteToFile writes text to a file with overwriting.
func WriteToFile(filePath string, data []byte) error {
	if config.Config.DryRun {
		config.Config.Out.Write(data)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		dir := filepath.Dir(filePath)
		os.MkdirAll(dir, 0700)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	return nil
}

// Copy copies src to dest, doesn't matter if src is a directory or a file.
func Copy(src, dest string) error {
	return copy.Copy(src, dest)
}
