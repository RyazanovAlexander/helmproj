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
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

// Run ...
func Run(projectFilePath string) error {
	project, _ := getProjectYaml(projectFilePath)

	/*result, err := project.Value([]string{"tree", "scalar"})
	if err != nil {
		return err
	}*/

	/*bytes, err := ioutil.ReadFile("C:\\SourceCode2\\helmproj\\examples\\charts\\backend\\values.yaml")
	if err != nil {
		return err
	}*/

	v := project.Values
	m := v["array"]
	c := m.([]interface{})
	ca := c[0].(string)
	fmt.Println(ca)

	m1 := v["map"]
	c1 := m1.([]interface{})
	ca1 := c1[0].(map[string]interface{})
	aa := ca1["prop11"]
	fmt.Println(aa)

	sourceYaml, _ := ioutil.ReadFile("C:\\SourceCode2\\helmproj\\examples\\charts\\backend\\values.yaml")
	var valuesYaml map[string]interface{}
	_ = yaml.Unmarshal(sourceYaml, &valuesYaml)
	result, _ := yaml.Marshal(valuesYaml)

	/*result, err = substituteMacroses(string(bytes), project)
	if err != nil {
		return err
	}*/

	tr := valuesYaml["tree"].(map[string]interface{})
	val, _ := project.Value2([]string{"array"})
	tr["array"] = val
	result, _ = yaml.Marshal(valuesYaml)

	fmt.Println(string(result))

	return nil
}

func getProjectYaml(projectFilePath string) (*Project, error) {
	sourceYaml, err := ioutil.ReadFile(projectFilePath)
	if err != nil {
		return nil, err
	}

	m := Project{}
	err = yaml.Unmarshal(sourceYaml, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func substituteMacroses(text string, project *Project) (string, error) {
	for {
		macrosPosition := ProjectRegexp.FindStringSubmatchIndex(text)
		if macrosPosition == nil {
			break
		}

		macros := text[macrosPosition[0]:macrosPosition[1]]
		firstReplacementPosition, err := getFirstReplacementPosition(text, macrosPosition[1])
		if err != nil {
			return "", err
		}

		macrosIdent := getMacrosIndent(text, macrosPosition[0])
		lastReplacementPosition := getLastReplacementPosition(text, firstReplacementPosition, macrosIdent)

		path := getPath(macros)
		projectValue, err := project.Value(path)
		projectValue = strings.TrimRightFunc(projectValue, endOfString)

		if err != nil {
			return "", err
		}

		text = replaceMacros(text, firstReplacementPosition, lastReplacementPosition, projectValue)
		fmt.Println(text)

		macrosPosition = ProjectRegexp.FindStringSubmatchIndex(text)
		text = removeLineWithMacros(text, macrosPosition[0])
		fmt.Println(text)
	}

	return text, nil
}

func getPath(rough string) []string {
	rough = strings.TrimRightFunc(rough, endOfString)
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

func getMacrosIndent(text string, position int) int {
	var ident = 0

	for i := position - 1; i > 0; i-- {
		if text[i] != ' ' {
			break
		}

		ident++
	}

	return ident
}

func getFirstReplacementPosition(text string, lastMacrosPosition int) (int, error) {
	for i := lastMacrosPosition; i < len(text); i++ {
		if text[i] == '#' {
			for j := i + 1; j < len(text); {
				if endOfString(rune(text[j])) {
					i = j + 1
					break
				}
				j++
				i = j
			}
		} else if text[i] == ':' {
			return i + 1, nil
		}
	}

	return 0, fmt.Errorf("The replaceable value described in the macro could not be found")
}

func getLastReplacementPosition(text string, startPosition int, ident int) int {
	lastReplacementPosition := startPosition

	var containsNewVariableOrMacros = func(text string, ident int) bool {
		macrosExists := ProjectRegexp.FindStringSubmatchIndex(text)
		if macrosExists != nil {
			return true
		}

		for i := 0; i < len(text) && ident >= 0; i++ {
			if text[i] == '#' {
				return false
			}
			if text[i] != ' ' {
				return true
			}
			ident--
		}
		return false
	}

	scanner := bufio.NewScanner(strings.NewReader(text[startPosition:]))
	var skipFirst = false
	newLine := 0
	for scanner.Scan() {
		line := scanner.Text()
		if !skipFirst {
			skipFirst = true
			lastReplacementPosition += len(line)
			continue
		}

		if containsNewVariableOrMacros(line, ident) {
			break
		}

		newLine++
		lastReplacementPosition += len(line)
	}

	if newLine > 1 {
		newLine--
	}

	return lastReplacementPosition + newLine
}

func replaceMacros(text string, firstReplacementPosition int, lastReplacementPosition int, projectValue string) string {
	return text[:firstReplacementPosition] + " " + projectValue + text[lastReplacementPosition:]
}

func removeLineWithMacros(text string, firstMacrosPosition int) string {
	result := ""
	scanner := bufio.NewScanner(strings.NewReader(text))

	cursor := 0
	for scanner.Scan() {
		line := scanner.Text()

		if cursor != -1 {
			cursor += len(line)
		}

		if cursor >= firstMacrosPosition {
			cursor = -1
			continue
		}

		result += fmt.Sprintln(line)
	}

	return result
}

func endOfString(c rune) bool {
	return c == '\r' || c == '\n'
}
