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
	"os"

	"github.com/otiai10/copy"
)

// DefaultValuesFile contains the name of a standard file with values.
const DefaultValuesFile string = "values.yaml"

// Chart describes the Chart
type Chart struct {
	values []Values
}

// LoadChart returns the model of the chart.
func LoadChart(chartPath string, additionValuesFiles []string) (*Chart, error) {
	chart := &Chart{}
	chart.loadValuesFiles(chartPath, additionValuesFiles)

	return chart, nil
}

// SubstituteValues substitutes values from the project file according to macros in values file.
func (chart *Chart) SubstituteValues(outputFolder string, tree map[string]interface{}) error {
	for _, values := range chart.values {
		err := values.SubstituteValues(tree)
		if err != nil {
			return err
		}

		values.SaveTo(outputFolder)
	}

	return nil
}

// SubstituteAppVersion substitutes appVersion from the project file.
func (chart *Chart) SubstituteAppVersion(appVersion string) error {

	return nil
}

// CopyTo copies the chart along the specified path.
func (chart *Chart) CopyTo(path string) error {
	/*if len(path) == 0 {
		return errors.New("Invalid parameter value passed: param - 'path', len(path) = 0")
	}

	bytes, err := yaml.Marshal(values)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}*/

	return nil
}

func (chart *Chart) loadValuesFiles(chartPath string, additionValuesFiles []string) error {
	valuesFiles := []string{chartPath + "/" + DefaultValuesFile}

	for _, filePath := range additionValuesFiles {
		valuesFiles = append(valuesFiles, chartPath+"/"+filePath)
	}

	for _, filePath := range valuesFiles {
		values := Values{}
		err := values.LoadValuesFile(filePath)

		if err != nil {
			return err
		}

		chart.values = append(chart.values, values)
	}

	return nil
}

func copyChart(fromFolderPath, toFolderPath string) error {
	err := os.RemoveAll(toFolderPath)
	if err != nil {
		return err
	}

	return copy.Copy(fromFolderPath, toFolderPath)
}
