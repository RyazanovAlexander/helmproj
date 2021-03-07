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

package project

import (
	"github.com/RyazanovAlexander/helmproj/v1/internal/io"
	"github.com/RyazanovAlexander/helmproj/v1/internal/yaml"
)

// Project describes the structure of the project file.
type Project struct {
	Values       map[string]interface{} `yaml:"values"`
	Charts       []ChartManifest        `yaml:"charts"`
	OutputFolder string                 `yaml:"outputFolder"`
}

// ChartManifest describes the structure of the chart.
type ChartManifest struct {
	Name            string   `yaml:"name"`
	Path            string   `yaml:"path"`
	AppVersion      string   `yaml:"appVersion"`
	AdditionlValues []string `yaml:"additionlValues"`
}

// LoadProjectFile returns the model of the project file.
func LoadProjectFile(filePath string) (*Project, error) {
	filePath, err := io.GetRuntimeFilePath(filePath)
	if err != nil {
		return nil, err
	}

	project := Project{}
	if err := yaml.UnmarshalFromFile(filePath, &project); err != nil {
		return nil, err
	}

	return &project, nil
}
