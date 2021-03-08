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

package cli

import (
	"os"
	"strconv"
	"testing"

	"github.com/spf13/cobra"

	"github.com/RyazanovAlexander/helmproj/v1/internal/test"
)

func TestNewSettings(t *testing.T) {
	defer test.ResetEnv()

	expected := true
	os.Setenv(HelmprojDebug, strconv.FormatBool(expected))
	settings := NewSettings()

	if expected != settings.Debug {
		t.Errorf("Debug = %v; want %v", settings.Debug, expected)
	}
}

func TestAddFlags(t *testing.T) {
	defer test.ResetEnv()

	expected := true
	os.Setenv(HelmprojDebug, strconv.FormatBool(expected))

	cmd := &cobra.Command{}
	flags := cmd.PersistentFlags()
	if flags.HasFlags() != false {
		t.Errorf("The flags structure is incorrectly initialized")
	}

	settings := NewSettings()
	settings.AddFlags(flags)

	if flags.HasFlags() != true {
		t.Errorf("Failed to retrieve flag")
	}
}
