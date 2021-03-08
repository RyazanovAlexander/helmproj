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
	"path/filepath"

	"github.com/RyazanovAlexander/helmproj/v1/internal/io"
	"github.com/RyazanovAlexander/helmproj/v1/internal/tree"
	"github.com/RyazanovAlexander/helmproj/v1/internal/yaml"
)

type values struct {
	name             string
	values           map[string]interface{}
	replaceablePairs []tree.ReplaceablePair
}

// Values describes the values.yaml file with replaceble pairs.
type Values interface {
	SubstituteValues(yamlTree map[string]interface{}) error
	SaveTo(chartPath, outputFolder string) error
}

// LoadValuesFile returns the model of the values file.
func LoadValuesFile(filePath string) (Values, error) {
	values := &values{}

	replaceablePairs, err := tree.GetReplaceablePairs(filePath)
	if err != nil {
		return nil, err
	}
	values.replaceablePairs = replaceablePairs

	values.name = filepath.Base(filePath)

	filePath, err = io.GetRuntimeFilePath(filePath)
	if err != nil {
		return nil, err
	}

	var valuesYaml map[string]interface{}
	if err := yaml.UnmarshalFromFile(filePath, &valuesYaml); err != nil {
		return nil, err
	}
	values.values = valuesYaml

	return values, nil
}

// SubstituteValues substitutes values from the project file according to macros in values file.
func (values *values) SubstituteValues(yamlTree map[string]interface{}) error {
	for _, replaceablePair := range values.replaceablePairs {
		projectValue, err := tree.GetNodeByPath(replaceablePair.ProjectPath, yamlTree)
		if err != nil {
			return err
		}

		if err = tree.SubstituteValueByPath(replaceablePair.ValuesPath, projectValue, values.values); err != nil {
			return err
		}
	}

	return nil
}

// SaveTo serializes values to the specified folder.
func (values *values) SaveTo(chartPath, outputFolder string) error {
	chartDir := filepath.Base(chartPath)

	filePath, err := io.GetRuntimeFilePath(filepath.Join(outputFolder, chartDir, values.name))
	if err != nil {
		return err
	}

	return yaml.MarshalToFile(filePath, values.values)
}
