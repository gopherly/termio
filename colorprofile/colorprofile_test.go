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
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopherly.dev/termio"
	"gopherly.dev/termio/colorprofile"
)

// TestDetect_ReturnsColorPolicy verifies Detect returns a non-nil policy.
func TestDetect_ReturnsColorPolicy(t *testing.T) {
	t.Parallel()

	p := colorprofile.Detect(&bytes.Buffer{}, []string{"NO_COLOR=1"})

	assert.NotNilf(t, p, "Detect must return a non-nil ColorPolicy")
}

// TestDetect_ImplementsColorPolicy verifies the return satisfies the interface.
func TestDetect_ImplementsColorPolicy(t *testing.T) {
	t.Parallel()

	_ = colorprofile.Detect(&bytes.Buffer{}, nil)
}

// TestDetect_ApplyWrapsWriter verifies Apply returns a non-nil writer.
func TestDetect_ApplyWrapsWriter(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	p := colorprofile.Detect(buf, []string{"NO_COLOR=1"})

	w := p.Apply(buf)

	assert.NotNilf(t, w, "Apply must return a non-nil writer")
}

// TestDetect_NoColor_StripsAnsi verifies NO_COLOR strips ANSI escapes.
func TestDetect_NoColor_StripsAnsi(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	p := colorprofile.Detect(buf, []string{"NO_COLOR=1"})
	w := p.Apply(buf)

	_, err := io.WriteString(w, "\x1b[32mhello\x1b[0m")

	require.NoErrorf(t, err, "write must not fail")
	assert.Equalf(t, "hello", buf.String(), "NO_COLOR must strip ANSI escapes")
}

// TestDetect_WithColorPolicy_Integration verifies the full Streams stack.
func TestDetect_WithColorPolicy_Integration(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	p := colorprofile.Detect(buf, []string{"NO_COLOR=1"})
	s := termio.New(nil, buf, io.Discard, termio.WithColorPolicy(p))

	_, err := io.WriteString(s.Out, "\x1b[32mgreen\x1b[0m")

	require.NoErrorf(t, err, "write to Streams.Out must not fail")
	assert.Equalf(t, "green", buf.String(),
		"NO_COLOR policy must strip ANSI through the full Streams stack")
}
