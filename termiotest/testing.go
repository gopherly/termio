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

// Package termiotest provides [*termio.Streams] test helpers with buffer-
// backed streams and convenience accessors for the underlying buffers.
//
// # Quick start
//
//	s, in, out, errOut := termiotest.New()
//	fmt.Fprintln(s.Out, "hello")
//	fmt.Println(out.String()) // prints "hello\n"
//
// Both [New] and [NewTTY] return the three [*bytes.Buffer] values
// corresponding to In, Out, and ErrOut so assertions are one line:
//
//	assert.Equal(t, "hello\n", out.String())
package termiotest

import (
	"bytes"

	"gopherly.dev/termio"
)

// New returns a [*termio.Streams] backed by three [*bytes.Buffer] values.
// All three streams report as non-TTY. The returned buffers are the live
// underlying storage; writes to s.Out appear in out, etc.
func New() (s *termio.Streams, in, out, errOut *bytes.Buffer) {
	in = &bytes.Buffer{}
	out = &bytes.Buffer{}
	errOut = &bytes.Buffer{}
	s = termio.New(in, out, errOut)

	return s, in, out, errOut
}

// NewTTY is like [New] but overrides all three TTY flags to true, making
// [termio.Streams.IsInteractive] return true. Use when the code path under
// test branches on TTY detection.
func NewTTY() (s *termio.Streams, in, out, errOut *bytes.Buffer) {
	s, in, out, errOut = New()
	s.SetStdinTTY(true)
	s.SetStdoutTTY(true)
	s.SetStderrTTY(true)

	return s, in, out, errOut
}
