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

package version

import (
	"runtime"
)

var (
	// Version is the current version of Helmproj.
	// Version is passed as an external parameter when building the program.
	// The version is of the format Major.Minor.Patch[-Prerelease].
	//
	// Increment major number for new feature additions and behavioral changes.
	// Increment minor number for bug fixes and performance enhancements.
	Version = "v0.0.0"

	// Buildtime is the build time of the program.
	Buildtime = ""

	// GitShortSHA is the git commit short sha1.
	GitShortSHA = ""
)

// BuildInfo describes the compile time information.
type BuildInfo struct {
	// Version is the current semver.
	Version string `json:"version,omitempty"`
	// GitShortSHA is the git commit short sha1.
	GitShortSHA string `json:"git_short_sha,omitempty"`
	// GitTreeState is the state of the git tree.
	Buildtime string `json:"build_time,omitempty"`
	// GoVersion is the version of the Go compiler used.
	GoVersion string `json:"go_version,omitempty"`
}

// GetBuildInfo returns build info
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:     Version,
		GitShortSHA: GitShortSHA,
		Buildtime:   Buildtime,
		GoVersion:   runtime.Version(),
	}
}
