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

import "io"

// Writer is a stream wrapper that layers three concerns on top of a raw
// [io.Writer]:
//
//   - Color adaptation: an optional [ColorPolicy] intercepts writes and
//     strips, translates, or passes through ANSI sequences.
//   - Sticky error: the first write error is latched and returned by all
//     subsequent writes, preventing silent partial output. The error is
//     independent per Writer, so an error on Out does not affect ErrOut.
//   - FD preservation: [Writer.Fd] returns the file descriptor of the original
//     underlying stream, letting downstream libraries (bubbletea, glamour,
//     lipgloss) detect the terminal even when they receive a wrapped value.
//
// Writer satisfies [io.Writer] and interface{ Fd() uintptr }.
// Obtain a Writer via [New] or [System]; do not construct one directly.
type Writer struct {
	w         io.Writer
	fd        uintptr
	stickyErr stickyError
}

// Write writes p to the underlying stream. If a previous write recorded an
// error, Write returns that error immediately without forwarding the bytes.
//
// Write implements [io.Writer].
func (w *Writer) Write(p []byte) (int, error) {
	if w.stickyErr.err != nil {
		return 0, w.stickyErr.err
	}

	n, err := w.w.Write(p)
	w.stickyErr.record(err)

	return n, err
}

// Fd returns the file descriptor of the underlying stream, or [InvalidFd]
// when the stream is not backed by a real OS file descriptor (for example,
// when a [bytes.Buffer] was supplied in tests).
//
// The returned value is a bare uintptr so that the Go garbage collector does
// not treat it as a live reference to the [os.File], matching the contract of
// [os.File.Fd]. Callers that need to call syscalls must ensure the File (or
// the Streams that owns this Writer) remains open for the duration of the
// operation.
func (w *Writer) Fd() uintptr { return w.fd }

// Err returns the first write error recorded by this Writer, or nil when all
// writes have succeeded so far. The error is latched: once set it is returned
// by every subsequent [Writer.Write] call and remains accessible via Err.
func (w *Writer) Err() error { return w.stickyErr.err }

// newWriter constructs a Writer over raw. When policy is non-nil, the write
// chain is: raw → policy.Apply(raw) → Writer. The FD is captured from raw
// before the policy wraps it, so Fd() always refers to the real descriptor.
func newWriter(raw io.Writer, policy ColorPolicy) *Writer {
	fd := fdOf(raw)

	adapted := raw
	if policy != nil {
		adapted = policy.Apply(raw)
	}

	return &Writer{w: adapted, fd: fd}
}

// stickyError latches the first non-nil error it is given.
type stickyError struct {
	err error
}

func (e *stickyError) record(err error) {
	if e.err == nil && err != nil {
		e.err = err
	}
}

// Compile-time interface checks.
var (
	_ io.Writer  = (*Writer)(nil)
	_ fileWriter = (*Writer)(nil)
)
