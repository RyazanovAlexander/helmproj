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
	"io"

	"github.com/spf13/cobra"

	"github.com/RyazanovAlexander/helmproj/v1/internal/cli"
	"github.com/RyazanovAlexander/helmproj/v1/internal/preprocessor"
)

var globalUsage = `This application is a tool for preprocessing values.yaml and Chart.yaml files.

Common actions for Helmproj:

- helmproj:                           preprocessing helm charts based on information from the ./project.yaml file
- helmproj -f ./path/to/project.yaml: preprocessing helm charts based on information from the project.yaml file located at the specified path

Environment variables:

| Name                               | Description                                                                       |
|------------------------------------|-----------------------------------------------------------------------------------|
| $HELMPROJ_DEBUG                    | indicate whether or not Helmprog is running in Debug mode                         |
`

// DefaultProjectFilePath is the path to the default project file.
const DefaultProjectFilePath = "./project.yaml"

// The path to the project file, passed as a command line argument.
var projFilePath string

// NewRootCmd creates new root cmd.
func NewRootCmd(out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "helmproj",
		Short: "Preprocessor for Helm Charts",
		Long:  globalUsage,
		Run:   func(cmd *cobra.Command, args []string) { runRootCmd(out, args) },
	}

	flags := cmd.PersistentFlags()
	flags.StringVarP(&projFilePath, "file", "f", DefaultProjectFilePath, "path to project file")
	flags.Parse(args)

	cli.Settings.AddFlags(flags)

	// Add subcommands
	cmd.AddCommand(
		newVersionCmd(out),
	)

	return cmd
}

func runRootCmd(out io.Writer, args []string) {
	//preprocessor.Run(projFilePath)
	//preprocessor.Run("C:\\SourceCode2\\helmproj\\examples\\charts\\project.yaml")
	preprocessor.Run("C:\\SourceCode2\\helmproj\\examples\\project.yaml") // REMOVE
}
