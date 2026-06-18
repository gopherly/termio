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

// Package termio provides terminal I/O primitives for Go CLI programs.
//
// # Overview
//
// [Streams] is the central type. It bundles an input reader (In) and two
// output writers (Out and ErrOut) together with TTY detection and terminal
// size querying. Call [System] to get streams backed by the real OS
// file descriptors:
//
//	s := termio.System()
//	fmt.Fprintln(s.Out, "hello, world")
//
// # Architecture
//
// The three streams are independent: Out and ErrOut are [*Writer] values
// rather than plain [io.Writer] values.  Each Writer layers three concerns
// on top of its raw stream:
//
//   - Color adaptation: an optional [ColorPolicy] intercepts writes and
//     strips, translates, or passes through ANSI sequences. When no policy
//     is configured, bytes reach the raw stream unmodified.
//   - Sticky errors: the first write error is latched per stream. An error
//     on Out does not affect ErrOut, and vice versa.
//   - FD preservation: [Writer.Fd] returns the original file descriptor,
//     so downstream libraries (bubbletea, glamour, lipgloss) can detect the
//     terminal even when they receive the wrapped value.
//
// The core package depends only on [golang.org/x/term]. Color adaptation
// is opt-in and lives in the [gopherly.dev/termio/colorprofile] sub-package;
// callers that do not need color never compile that dependency.
//
// # Quick start - with color adaptation
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
// # Quick start - without color
//
//	s := termio.System()  // no color dep compiled
//	fmt.Fprintln(s.Out, "plain output")
//
// # Testing
//
// The [gopherly.dev/termio/termiotest] package provides buffer-backed
// helpers that return the underlying [*bytes.Buffer] values for assertion:
//
//	s, _, out, _ := termiotest.New()
//	fmt.Fprintln(s.Out, "test output")
//	assert.Equal(t, "test output\n", out.String())
//
// # Exported surface
//
//	Type          Description
//	Streams       central I/O bundle (In, Out, ErrOut)
//	Writer        stream wrapper (color + sticky error + FD)
//	ColorPolicy   interface for ANSI color adaptation
//	Option        functional option for New / System
//
//	Constructor   Description
//	System        streams from os.Stdin / os.Stdout / os.Stderr
//	New           streams from caller-supplied readers/writers
//
//	Option        Description
//	WithColorPolicy  inject a ColorPolicy into Out and ErrOut
//
//	Constant       Description
//	DefaultWidth   fallback column width (80) when width is unknown
//	DefaultHeight  fallback row height (24) when height is unknown
//	InvalidFd      sentinel Fd value (^uintptr(0)) for non-file streams
package termio
