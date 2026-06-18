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
//
// # Constructors
//
//	Constructor   Description
//	Detect        detect profile from an io.Writer and environment
//	From          build a Policy from an explicit Profile level
//
// # Exported surface
//
//	Type      Description
//	Policy    ColorPolicy backed by a detected or explicit Profile level
//	Profile   alias for charmbracelet/colorprofile.Profile
//
//	Constant   Description
//	NoTTY      no terminal attached
//	ASCII      no color support
//	ANSI       16-color ANSI
//	ANSI256    256-color ANSI
//	TrueColor  24-bit true color
package colorprofile
