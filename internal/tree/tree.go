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

package tree

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/karalabe/cookiejar.v2/collections/stack"

	"github.com/RyazanovAlexander/helmproj/v1/internal/io"
)

// ReplaceablePair describes the structure of paths to replaceable pairs.
type ReplaceablePair struct {
	ValuesPath  []string
	ProjectPath []string
}

// ProjectMacrosRegexp is a regular expression to search for a macro with a project.
var ProjectMacrosRegexp = regexp.MustCompile("# *{{ *.Project.*}}(\r\n|\r|\n|)")

// CommentRegexp is a regex to find comments.
var CommentRegexp = regexp.MustCompile("^ *#")

// VariableRegexp is a regex to find variables.
var VariableRegexp = regexp.MustCompile(".+:")

// GetReplaceablePairs returns replaceable pairs.
//
// Example:
// -----------------------------
// scalar: value {{ .Project.scalar }}
//
// tree:
//   scalar: value # {{ .Project.tree.value }}
//   map: # {{ .Project.someval }}
//     - key1: val1
//       prop1: prop-val1
// -----------------------------
//
// Result:
//             values.yaml             project.yaml
// result[0] = { "scalar" }         -> { "scalar" }
// result[1] = { "tree", "scalar" } -> { "tree", "value" }
// result[2] = { "tree", "map" }    -> { "someval" }
func GetReplaceablePairs(valuesFilePath string) ([]ReplaceablePair, error) {
	var replaceablePairs []ReplaceablePair

	valuesFilePath, err := io.GetRuntimeFilePath(valuesFilePath)
	if err != nil {
		return nil, err
	}

	valuesFile, err := os.Open(valuesFilePath)
	if err != nil {
		return nil, err
	}
	defer valuesFile.Close()

	variableTreeStack := stack.New()

	scanner := bufio.NewScanner(valuesFile)
	for scanner.Scan() {
		line := scanner.Text()

		var projectPath []string

		projectMacros := ProjectMacrosRegexp.FindStringSubmatch(line)
		comment := CommentRegexp.FindStringSubmatch(line)
		if projectMacros != nil {
			projectPath = getPath(projectMacros[0])
		} else if comment != nil {
			continue
		}

		variable := VariableRegexp.FindStringSubmatch(line)
		if variable != nil {
			variable := strings.Split(variable[0], ":")[0]
			variableTreeStack.Push(variable)
		}

		if projectPath != nil {
			var valuesPath []string
			var tmp []string

			for {
				blankVariable := variableTreeStack.Top().(string)
				variable := strings.TrimLeft(blankVariable, " ")
				valuesPath = append(valuesPath, variable)

				if !haveLeadingSpace(blankVariable) {
					break
				}

				tmp = append(tmp, variableTreeStack.Pop().(string))
			}

			for i := len(tmp) - 1; i >= 1; i-- {
				variableTreeStack.Push(tmp[i])
			}

			valuesPath = reverse(valuesPath)

			replaceablePairs = append(replaceablePairs, ReplaceablePair{
				ValuesPath:  valuesPath,
				ProjectPath: projectPath,
			})
		}
	}

	return replaceablePairs, nil
}

// GetNodeByPath returns a node at the specified path.
func GetNodeByPath(path []string, node interface{}) (interface{}, error) {
	if len(path) == 0 {
		return "", errors.New("Invalid parameter value passed: param - 'path', len(path) = 0")
	}

	var runner interface{} = node

	for i := 0; i < len(path); i++ {
		value, ok := runner.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("No value found in the specified path %v", path)
		}

		runner = value[path[i]]
	}

	return runner, nil
}

// SubstituteValueByPath replaces the value at the specified path.
func SubstituteValueByPath(path []string, node interface{}, tree interface{}) error {
	if len(path) == 0 {
		return errors.New("Invalid parameter value passed: param - 'path', len(path) = 0")
	}

	var runner interface{} = tree

	for i := 0; i < len(path); i++ {
		value, ok := runner.(map[string]interface{})
		if !ok {
			return fmt.Errorf("No value found in the specified path %v", path)
		}

		if i == len(path)-1 {
			value[path[i]] = node
			break
		}

		runner = value[path[i]]
	}

	return nil
}

func haveLeadingSpace(line string) bool {
	return line[0] == ' '
}

func getPath(rough string) []string {
	rough = strings.TrimLeft(rough, "#")
	rough = strings.TrimSpace(rough)
	rough = strings.TrimLeft(rough, "{{")
	rough = strings.TrimRight(rough, "}}")
	rough = strings.TrimSpace(rough)
	rough = strings.TrimLeft(rough, ".")
	rough = strings.TrimLeft(rough, "Project")
	rough = strings.TrimLeft(rough, ".")

	return strings.Split(rough, ".")
}

func reverse(values []string) []string {
	for i := 0; i < len(values)/2; i++ {
		j := len(values) - i - 1
		values[i], values[j] = values[j], values[i]
	}
	return values
}
