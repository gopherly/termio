// Copyright 2026 The Gopherly Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package colorprofile provides a [termio.ColorPolicy] implementation backed
// by [github.com/charmbracelet/colorprofile].
//
// This is the only package in the termio module that depends on
// charmbracelet/colorprofile. Programs that import only [gopherly.dev/termio]
// never compile this dependency.
//
// # Quick start
//
//	import (
//	    "os"
//	    "gopherly.dev/termio"
//	    "gopherly.dev/termio/colorprofile"
//	)
//
//	s := termio.System(
//	    termio.WithColorPolicy(colorprofile.Detect(os.Stdout, os.Environ())),
//	)
//	fmt.Fprintln(s.Out, "\x1b[32mgreen\x1b[0m or plain, depending on the terminal")
package colorprofile

import (
	"io"

	"gopherly.dev/termio"

	cbcolorprofile "github.com/charmbracelet/colorprofile"
)

// Detect returns a [termio.ColorPolicy] that adapts ANSI escape sequences
// based on the color capability of output, consulting env for environment
// variables such as NO_COLOR, COLORTERM, and TERM.
//
// When output is an *[os.File] the profile is detected against the real
// file descriptor; when it is any other writer (e.g. a [*bytes.Buffer] in
// tests), the detection falls back to the environment alone.
//
// The returned policy is safe for concurrent use.
func Detect(output io.Writer, env []string) termio.ColorPolicy { //nolint:ireturn
	profile := cbcolorprofile.Detect(output, env)

	return &profilePolicy{profile: profile}
}

// profilePolicy wraps a charmbracelet/colorprofile Profile as a
// [termio.ColorPolicy].
type profilePolicy struct {
	profile cbcolorprofile.Profile
}

// Apply wraps w with a charmbracelet/colorprofile.Writer that rewrites ANSI
// escape sequences according to the detected terminal capability.
func (p *profilePolicy) Apply(w io.Writer) io.Writer {
	return &cbcolorprofile.Writer{Forward: w, Profile: p.profile}
}
