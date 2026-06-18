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

package colorprofile

import (
	"io"

	"gopherly.dev/termio"

	cbcolorprofile "github.com/charmbracelet/colorprofile"
)

// Profile is the detected color capability level of a terminal.
// It is an alias for [github.com/charmbracelet/colorprofile.Profile].
type Profile = cbcolorprofile.Profile

// Color capability levels, re-exported from charmbracelet/colorprofile for
// use without a second import.
const (
	NoTTY     = cbcolorprofile.NoTTY
	ASCII     = cbcolorprofile.ASCII
	ANSI      = cbcolorprofile.ANSI
	ANSI256   = cbcolorprofile.ANSI256
	TrueColor = cbcolorprofile.TrueColor
)

// Policy adapts ANSI color output to a terminal's detected capability level.
// Use [Detect] for auto-detection or [From] to set an explicit level.
// Policy implements [termio.ColorPolicy].
type Policy struct {
	profile Profile
}

// Detect returns a [Policy] that adapts ANSI escape sequences based on the
// color capability of output, consulting env for environment variables such
// as NO_COLOR, COLORTERM, and TERM.
//
// When output is an *[os.File] the profile is detected against the real file
// descriptor; when it is any other writer (e.g. a [*bytes.Buffer] in tests),
// the detection falls back to the environment alone.
//
// The returned policy is safe for concurrent use.
func Detect(output io.Writer, env []string) *Policy {
	return &Policy{profile: cbcolorprofile.Detect(output, env)}
}

// From returns a [Policy] for an explicit profile level, bypassing detection.
// Use this when the level is known ahead of time (e.g. from a --color flag).
func From(p Profile) *Policy {
	return &Policy{profile: p}
}

// Profile returns the color capability level this policy enforces.
func (p *Policy) Profile() Profile { return p.profile }

// String returns the human-readable name of the profile level (e.g. "TrueColor").
func (p *Policy) String() string { return p.profile.String() }

// Apply wraps w with a charmbracelet/colorprofile.Writer that rewrites ANSI
// escape sequences according to the policy profile level.
func (p *Policy) Apply(w io.Writer) io.Writer {
	return &cbcolorprofile.Writer{Forward: w, Profile: p.profile}
}

// Compile-time interface check.
var _ termio.ColorPolicy = (*Policy)(nil)
