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

package chart

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/RyazanovAlexander/helmproj/v1/internal/io"
)

// DefaultValuesFile contains the name of a standard file with values.
const DefaultValuesFile string = "values.yaml"

// ChartFileName is a default chart file name.
const ChartFileName string = "Chart.yaml"

// AppVersionTemplate describes the template for the variable appversion in Chart.yaml file.
const AppVersionTemplate string = "appVersion: %s"

// ProjectMacrosRegexp is a regular expression to search for a chart version in Chart.yaml file.
var ProjectMacrosRegexp = regexp.MustCompile(".*appVersion:.*(\r\n|\r|\n|)")

type chart struct {
	path       string
	appVersion string
	values     []Values
}

// Chart describes the Chart and defines the available operations in relation to it
type Chart interface {
	SubstituteValues(outputFolder string, tree map[string]interface{}) error
	SubstituteAppVersion(outputFolder, appVersion string) error
	CopyTo(toPath string) error
}

// LoadChart returns the model of the chart.
func LoadChart(chartPath, appVersion string, additionValuesFiles []string) (Chart, error) {
	chart := &chart{}
	chart.path = chartPath
	chart.appVersion = appVersion
	chart.loadValuesFiles(chartPath, additionValuesFiles)

	return chart, nil
}

// SubstituteValues substitutes values from the project file according to macros in values file.
func (chart *chart) SubstituteValues(outputFolder string, tree map[string]interface{}) error {
	for _, values := range chart.values {
		err := values.SubstituteValues(tree)
		if err != nil {
			return err
		}

		values.SaveTo(chart.path, outputFolder)
	}

	return nil
}

// SubstituteAppVersion substitutes appVersion from the project file.
func (chart *chart) SubstituteAppVersion(outputFolder, appVersion string) error {
	chartDir := filepath.Base(chart.path)
	chartFilePath, err := io.GetRuntimeFilePath(filepath.Join(outputFolder, chartDir, ChartFileName))
	if err != nil {
		return err
	}

	fileContent, err := ioutil.ReadFile(chartFilePath)
	if err != nil {
		return err
	}

	location := ProjectMacrosRegexp.FindIndex(fileContent)
	if location == nil {
		return errors.New("appVersion section not found in Chart.yaml file")
	}

	oldString := fileContent[location[0]:location[1]]
	newString := fmt.Sprintf(AppVersionTemplate, chart.appVersion)

	newContent := strings.Replace(string(fileContent), string(oldString), newString, 1)
	return io.WriteToFile(chartFilePath, []byte(newContent))
}

// CopyTo copies the chart along the specified path.
func (chart *chart) CopyTo(toPath string) error {
	if len(toPath) == 0 {
		return errors.New("Invalid parameter value passed: param - 'path', len(path) = 0")
	}

	toPath, err := io.GetRuntimeFilePath(toPath)
	if err != nil {
		return err
	}

	src, err := io.GetRuntimeFilePath(chart.path)
	if err != nil {
		return err
	}

	dest := filepath.Join(toPath, filepath.Base(chart.path))
	return io.Copy(src, dest)
}

func (chart *chart) loadValuesFiles(chartPath string, additionValuesFiles []string) error {
	valuesFiles := []string{filepath.Join(chartPath, DefaultValuesFile)}

	for _, filePath := range additionValuesFiles {
		valuesFiles = append(valuesFiles, filepath.Join(chartPath, filePath))
	}

	for _, filePath := range valuesFiles {
		values, err := LoadValuesFile(filePath)

		if err != nil {
			return err
		}

		chart.values = append(chart.values, values)
	}

	return nil
}
