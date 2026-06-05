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

package termio

import (
	"errors"
	"io"
	"math"
	"os"

	"golang.org/x/term"
)

// DefaultWidth is the terminal column width assumed when the underlying
// stream is not a terminal or its size cannot be determined.
const DefaultWidth = 80

// Streams bundles the three standard I/O channels a CLI program needs —
// input, output, and diagnostics — together with terminal capability
// detection. It is the central type of the termio package.
//
// Out and ErrOut are [*Writer] values rather than plain [io.Writer]
// values: the concrete type surfaces per-stream sticky-error state and FD
// preservation without type assertions. Color adaptation is injected via
// [WithColorPolicy]; when no policy is configured the raw stream is used
// directly. An error on Out does not affect ErrOut, and vice versa.
//
// Streams is safe to construct concurrently with other operations but is
// not safe for concurrent mutation via the SetXxxTTY methods.
type Streams struct {
	// In is the input stream. It is the user-supplied [io.Reader] (typically
	// [os.Stdin]) and is not wrapped. Treat as read-only after construction.
	In io.Reader

	// Out is the primary output stream. Writes are color-adapted (when a
	// [ColorPolicy] is set), the first error is latched, and the original
	// FD is preserved for terminal-detection by downstream libraries.
	Out *Writer

	// ErrOut is the diagnostics stream. It has the same wrapping policy as
	// Out, but its sticky error is independent.
	ErrOut *Writer

	rawIn     io.Reader
	rawOut    io.Writer
	rawErrOut io.Writer

	stdinIsTTY  bool
	stdoutIsTTY bool
	stderrIsTTY bool
}

// System returns a [Streams] backed by [os.Stdin], [os.Stdout], and
// [os.Stderr]. TTY status and terminal width are detected against the real
// file descriptors. Apply functional options to configure color adaptation.
func System(opts ...Option) *Streams {
	return New(os.Stdin, os.Stdout, os.Stderr, opts...)
}

// New returns a [Streams] over the supplied streams. Pass *[os.File] values
// for production code (TTY detection works against the real file descriptor).
// For tests, prefer [gopherly.dev/termio/termiotest.New], which also returns
// the underlying buffers for assertion.
//
// Nil arguments are replaced with safe no-op streams: a reader that returns
// [io.EOF] immediately and a writer that discards all bytes.
func New(in io.Reader, out, errOut io.Writer, opts ...Option) *Streams {
	if in == nil {
		in = nopReader{}
	}

	if out == nil {
		out = io.Discard
	}

	if errOut == nil {
		errOut = io.Discard
	}

	cfg := &config{}
	applyAll(cfg, opts)

	return &Streams{
		In:          in,
		Out:         newWriter(out, cfg.colorPolicy),
		ErrOut:      newWriter(errOut, cfg.colorPolicy),
		rawIn:       in,
		rawOut:      out,
		rawErrOut:   errOut,
		stdinIsTTY:  isTerminal(in),
		stdoutIsTTY: isTerminal(out),
		stderrIsTTY: isTerminal(errOut),
	}
}

// IsInteractive reports whether both stdin and stdout are terminals. When
// true, the program may safely display interactive prompts.
func (s *Streams) IsInteractive() bool {
	return s.stdinIsTTY && s.stdoutIsTTY
}

// IsStdinTTY reports whether the input stream is a terminal.
func (s *Streams) IsStdinTTY() bool { return s.stdinIsTTY }

// IsStdoutTTY reports whether the primary output stream is a terminal.
func (s *Streams) IsStdoutTTY() bool { return s.stdoutIsTTY }

// IsStderrTTY reports whether the diagnostics stream is a terminal.
func (s *Streams) IsStderrTTY() bool { return s.stderrIsTTY }

// SetStdinTTY overrides the cached stdin TTY status. Use in tests to
// test interactive code paths with a buffer-backed stream.
func (s *Streams) SetStdinTTY(v bool) { s.stdinIsTTY = v }

// SetStdoutTTY overrides the cached stdout TTY status. See
// [Streams.SetStdinTTY] for typical usage.
func (s *Streams) SetStdoutTTY(v bool) { s.stdoutIsTTY = v }

// SetStderrTTY overrides the cached stderr TTY status. See
// [Streams.SetStdinTTY] for typical usage.
func (s *Streams) SetStderrTTY(v bool) { s.stderrIsTTY = v }

// TerminalWidth returns the column width of the controlling terminal, or
// [DefaultWidth] when stdout is not a terminal or its size cannot be queried.
func (s *Streams) TerminalWidth() int {
	fd := s.Out.Fd()
	if fd == InvalidFd {
		return DefaultWidth
	}

	if fd > uintptr(math.MaxInt) {
		return DefaultWidth
	}

	w, _, err := term.GetSize(int(fd))
	if err != nil || w <= 0 {
		return DefaultWidth
	}

	return w
}

// Err returns the first write error encountered by either Out or ErrOut, or
// nil when all writes have succeeded so far. When both streams have latched
// errors, both are reported joined with [errors.Join].
func (s *Streams) Err() error {
	return errors.Join(s.Out.Err(), s.ErrOut.Err())
}

// RawIn returns the unwrapped input stream supplied to [New] or [System].
func (s *Streams) RawIn() io.Reader { return s.rawIn }

// RawOut returns the unwrapped output stream supplied to [New] or [System].
func (s *Streams) RawOut() io.Writer { return s.rawOut }

// RawErrOut returns the unwrapped diagnostics stream supplied to [New] or
// [System].
func (s *Streams) RawErrOut() io.Writer { return s.rawErrOut }

// isTerminal reports whether v is backed by a real terminal.
func isTerminal(v any) bool {
	fd := fdOf(v)
	if fd == InvalidFd {
		return false
	}

	if fd > uintptr(math.MaxInt) {
		return false
	}

	return term.IsTerminal(int(fd))
}

// nopReader is a reader that returns [io.EOF] on every Read call.
type nopReader struct{}

func (nopReader) Read(_ []byte) (int, error) { return 0, io.EOF }
