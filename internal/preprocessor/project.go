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

package preprocessor

import (
	"errors"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

// ProjectRegexp is a regular expression to search for a macro with a project
var ProjectRegexp = regexp.MustCompile("# *{{ .Project.*}}(\r\n|\r|\n|)")

// Project describes the structure of the project file.
type Project struct {
	Values map[string]interface{} `yaml:"values"`
	Charts []struct {
		Name       string `yaml:"name"`
		Path       string `yaml:"path"`
		AppVersion string `yaml:"appVersion"`
	} `yaml:"charts"`
	OutputFolder string `yaml:"outputFolder"`
}

// Value returns a value at the specified path.
func (project *Project) Value(path []string) (string, error) {
	if len(path) == 0 {
		return "", errors.New("Invalid parameter value passed: param - 'path', len(path) = 0")
	}

	var runner interface{} = project.Values[path[0]]

	for i := 1; i < len(path); i++ {
		if i == len(path) {
			if runner == nil {
				return "", fmt.Errorf("No value found in the specified path %v", path)
			}
			break
		}

		value, ok := runner.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("No value found in the specified path %v", path)
		}

		runner = value[path[1]]
	}

	result, err := yaml.Marshal(&runner)

	return string(result), err
}

// Value returns a value at the specified path.
func (project *Project) Value2(path []string) (interface{}, error) {
	if len(path) == 0 {
		return "", errors.New("Invalid parameter value passed: param - 'path', len(path) = 0")
	}

	var runner interface{} = project.Values[path[0]]

	for i := 1; i < len(path); i++ {
		if i == len(path) {
			if runner == nil {
				return "", fmt.Errorf("No value found in the specified path %v", path)
			}
			break
		}

		value, ok := runner.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("No value found in the specified path %v", path)
		}

		runner = value[path[1]]
	}

	return runner, nil
}
