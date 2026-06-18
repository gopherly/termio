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

package colorprofile_test

import (
	"bytes"
	"fmt"
	"io"

	"gopherly.dev/termio"
	"gopherly.dev/termio/colorprofile"
)

// ExamplePolicy_Profile shows how to build a policy from an explicit level
// and retrieve it. This is useful when passing the level to Charm libraries
// (lipgloss, glamour) or branching on color capability.
func ExamplePolicy_Profile() {
	policy := colorprofile.From(colorprofile.ANSI256)
	fmt.Println(policy.Profile())
	// Output:
	// ANSI256
}

// ExampleDetect shows how to wire colorprofile.Detect into termio.System.
// When NO_COLOR is set, ANSI escape sequences are stripped before reaching
// the output buffer.
func ExampleDetect() {
	buf := &bytes.Buffer{}
	env := []string{"NO_COLOR=1"}

	p := colorprofile.Detect(buf, env)
	s := termio.New(nil, buf, io.Discard, termio.WithColorPolicy(p))

	fmt.Fprintln(s.Out, "\x1b[32mhello\x1b[0m") //nolint:errcheck
	fmt.Print(buf.String())
	// Output:
	// hello
}

// ExampleFrom builds a Policy from an explicit profile level, bypassing
// auto-detection. Use this when the level is known ahead of time.
func ExampleFrom() {
	policy := colorprofile.From(colorprofile.ANSI)
	fmt.Println(policy.Profile())
	// Output:
	// ANSI
}
