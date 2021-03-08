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

package yaml

import (
	"io/ioutil"

	"github.com/RyazanovAlexander/helmproj/v1/internal/io"
	"gopkg.in/yaml.v3"
)

// UnmarshalFromFile deserializes the yaml file to the specified data type.
func UnmarshalFromFile(filePath string, object interface{}) error {
	sourceYaml, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(sourceYaml, object)
	if err != nil {
		return err
	}

	return nil
}

// MarshalToFile serialises the specified data type to the yaml file.
func MarshalToFile(filePath string, object interface{}) error {
	bytes, err := yaml.Marshal(object)
	if err != nil {
		return err
	}

	io.WriteToFile(filePath, bytes)

	return nil
}
