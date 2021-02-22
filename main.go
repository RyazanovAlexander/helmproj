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

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RyazanovAlexander/helmproj/v1/cmd"
	"github.com/RyazanovAlexander/helmproj/v1/internal/flog"
)

var Version string
var Buildtime string

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	os.Stdout.WriteString(fmt.Sprintf("Version: %s\n", Version))
	os.Stdout.WriteString(fmt.Sprintf("Buildtime: %s\n", Buildtime))

	cmd, err := cmd.NewRootCmd(os.Stdout, os.Args[1:])
	if err != nil {
		flog.WarningF("%+v", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		flog.DebugF("%+v", err)
		os.Exit(1)
	}
}
