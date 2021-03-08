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

package cmd

import (
	"bytes"
	"os"
	"testing"

	shellwords "github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"

	"github.com/RyazanovAlexander/helmproj/v1/internal/test"
)

// TestCase describes a test case that works with releases.
type TestCase struct {
	Name      string
	Cmd       string
	Golden    string
	WantError bool
	Repeat    int
}

// RunTestCmd runs the specified test cases
func RunTestCmd(t *testing.T, tests []TestCase) {
	t.Helper()

	for _, testCase := range tests {
		for i := 0; i <= testCase.Repeat; i++ {
			t.Run(testCase.Name, func(t *testing.T) {
				defer test.ResetEnv()()

				t.Logf("running cmd (attempt %d): %s", i+1, testCase.Cmd)
				_, out, err := executeActionCommandC(testCase.Cmd)

				if (err != nil) != testCase.WantError {
					t.Errorf("expected error, got '%v'", err)
				}

				if testCase.Golden != "" {
					test.AssertGoldenString(t, out, testCase.Golden)
				}
			})
		}
	}
}

func executeActionCommandC(cmd string) (*cobra.Command, string, error) {
	args, err := shellwords.Parse(cmd)
	if err != nil {
		return nil, "", err
	}

	buf := new(bytes.Buffer)
	rootCmd := NewRootCmd(buf, args)

	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	oldStdin := os.Stdin

	command, err := rootCmd.ExecuteC()
	result := buf.String()
	os.Stdin = oldStdin

	return command, result, err
}
