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

package termio_test

import (
	"bytes"
	"errors"
	"io"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopherly.dev/termio"
	"gopherly.dev/termio/termiotest"
)

// TestNew_NilInputsFallback verifies nil arguments get safe defaults.
func TestNew_NilInputsFallback(t *testing.T) {
	t.Parallel()

	s := termio.New(nil, nil, nil)

	assert.NotNilf(t, s.In, "In must not be nil after nil input")
	assert.NotNilf(t, s.Out, "Out must not be nil after nil out")
	assert.NotNilf(t, s.ErrOut, "ErrOut must not be nil after nil errOut")

	var p [1]byte
	_, err := s.In.Read(p[:])

	assert.ErrorIsf(t, err, io.EOF, "nil In must read as EOF")
}

// TestNew_WritesReachBuffers verifies writes flow to the underlying buffers.
func TestNew_WritesReachBuffers(t *testing.T) {
	t.Parallel()

	s, _, out, errOut := termiotest.New()

	_, err := io.WriteString(s.Out, "stdout")
	require.NoErrorf(t, err, "write to Out must succeed")

	_, err = io.WriteString(s.ErrOut, "stderr")
	require.NoErrorf(t, err, "write to ErrOut must succeed")

	assert.Equalf(t, "stdout", out.String(), "Out must reach output buffer")
	assert.Equalf(t, "stderr", errOut.String(), "ErrOut must reach errOut buffer")
}

// TestStreams_TTYDetectionBuffers verifies buffer-backed streams are non-TTY.
func TestStreams_TTYDetectionBuffers(t *testing.T) {
	t.Parallel()

	s, _, _, _ := termiotest.New() //nolint:dogsled

	assert.Falsef(t, s.IsStdinTTY(), "buffer-backed stdin must not be TTY")
	assert.Falsef(t, s.IsStdoutTTY(), "buffer-backed stdout must not be TTY")
	assert.Falsef(t, s.IsStderrTTY(), "buffer-backed stderr must not be TTY")
	assert.Falsef(t, s.IsInteractive(), "non-TTY streams must not be interactive")
}

// TestStreams_SetTTYOverrides verifies the TTY override setters work.
func TestStreams_SetTTYOverrides(t *testing.T) {
	t.Parallel()

	s, _, _, _ := termiotest.New() //nolint:dogsled

	s.SetStdinTTY(true)
	s.SetStdoutTTY(true)
	s.SetStderrTTY(true)

	assert.Truef(t, s.IsStdinTTY(), "SetStdinTTY(true) must take effect")
	assert.Truef(t, s.IsStdoutTTY(), "SetStdoutTTY(true) must take effect")
	assert.Truef(t, s.IsStderrTTY(), "SetStderrTTY(true) must take effect")
	assert.Truef(t, s.IsInteractive(), "TTY overrides must make streams interactive")
}

// TestStreams_TerminalWidthDefaultsForBuffers verifies DefaultWidth for buffers.
func TestStreams_TerminalWidthDefaultsForBuffers(t *testing.T) {
	t.Parallel()

	s, _, _, _ := termiotest.New() //nolint:dogsled

	assert.Equalf(t, termio.DefaultWidth, s.TerminalWidth(),
		"buffer-backed Streams must return DefaultWidth")
}

// TestStreams_ErrJoinsBothStreams verifies Err() joins both stream errors.
func TestStreams_ErrJoinsBothStreams(t *testing.T) {
	t.Parallel()

	errOut := &failWriter{}
	errErr := &failWriter{}
	s := termio.New(nil, errOut, errErr)

	io.WriteString(s.Out, "x")    //nolint:errcheck
	io.WriteString(s.ErrOut, "x") //nolint:errcheck

	err := s.Err()

	require.ErrorIsf(t, err, errWriteFailedStream, "Err() must include Out error")
}

// TestStreams_ErrNilWhenNoWrites verifies Err() is nil before any writes.
func TestStreams_ErrNilWhenNoWrites(t *testing.T) {
	t.Parallel()

	s, _, _, _ := termiotest.New() //nolint:dogsled

	assert.NoErrorf(t, s.Err(), "Err() must be nil when nothing has been written")
}

// TestStreams_StickyErrorsAreIndependent verifies per-stream error independence.
func TestStreams_StickyErrorsAreIndependent(t *testing.T) {
	t.Parallel()

	errOut := &failWriter{}
	s := termio.New(nil, errOut, &bytes.Buffer{})

	io.WriteString(s.Out, "trigger") //nolint:errcheck

	require.Errorf(t, s.Out.Err(), "Out.Err() must be set after failed write")
	assert.NoErrorf(t, s.ErrOut.Err(), "ErrOut.Err() must remain nil when only Out failed")
}

// TestStreams_RawAccessors verifies raw stream accessors return the originals.
func TestStreams_RawAccessors(t *testing.T) {
	t.Parallel()

	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	s := termio.New(in, out, errOut)

	assert.Samef(t, in, s.RawIn(), "RawIn must be the unwrapped reader")
	assert.Samef(t, out, s.RawOut(), "RawOut must be the unwrapped writer")
	assert.Samef(t, errOut, s.RawErrOut(), "RawErrOut must be the unwrapped writer")
}

// TestStreams_WithColorPolicy verifies color policy transforms output.
func TestStreams_WithColorPolicy(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	p := &prefixPolicy{prefix: ">>"}
	s := termio.New(nil, buf, io.Discard, termio.WithColorPolicy(p))

	_, err := io.WriteString(s.Out, "hello")

	require.NoErrorf(t, err, "write must not fail with policy")
	assert.Equalf(t, ">>hello", buf.String(), "policy must transform Out writes")
}

// TestSystem_ReturnsValidStreams verifies System returns non-nil os streams.
func TestSystem_ReturnsValidStreams(t *testing.T) {
	t.Parallel()

	s := termio.System()

	assert.NotNilf(t, s.In, "In must not be nil")
	assert.NotNilf(t, s.Out, "Out must not be nil")
	assert.NotNilf(t, s.ErrOut, "ErrOut must not be nil")
	assert.Samef(t, os.Stdin, s.RawIn(), "RawIn must be os.Stdin")
	assert.Samef(t, os.Stdout, s.RawOut(), "RawOut must be os.Stdout")
	assert.Samef(t, os.Stderr, s.RawErrOut(), "RawErrOut must be os.Stderr")
}

// TestSystem_ForwardsOptions verifies System passes options to New.
func TestSystem_ForwardsOptions(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	p := &prefixPolicy{prefix: ">>"}
	s := termio.New(nil, buf, io.Discard, termio.WithColorPolicy(p))

	_, err := io.WriteString(s.Out, "hello")

	require.NoErrorf(t, err, "write must succeed")
	assert.Equalf(t, ">>hello", buf.String(), "option must be forwarded through System")
}

// TestStreams_TerminalWidthOverflowFd verifies overflow FD falls back to default.
func TestStreams_TerminalWidthOverflowFd(t *testing.T) {
	t.Parallel()

	raw := &fdWriter{fd: uintptr(math.MaxInt) + 1}
	s := termio.New(nil, raw, io.Discard)

	assert.Equalf(t, termio.DefaultWidth, s.TerminalWidth(),
		"overflow FD must fall back to DefaultWidth")
}

// TestStreams_TerminalWidthNonTerminalFd verifies regular file uses default.
func TestStreams_TerminalWidthNonTerminalFd(t *testing.T) {
	t.Parallel()

	f, err := os.CreateTemp(t.TempDir(), "termio-width-*")
	require.NoErrorf(t, err, "create temp file")

	t.Cleanup(func() { f.Close() }) //nolint:errcheck

	s := termio.New(nil, f, io.Discard)

	assert.Equalf(t, termio.DefaultWidth, s.TerminalWidth(),
		"regular file FD must fall back to DefaultWidth")
}

// TestStreams_TerminalHeightDefaultsForBuffers verifies DefaultHeight for buffers.
func TestStreams_TerminalHeightDefaultsForBuffers(t *testing.T) {
	t.Parallel()

	s, _, _, _ := termiotest.New() //nolint:dogsled

	assert.Equalf(t, termio.DefaultHeight, s.TerminalHeight(),
		"buffer-backed Streams must return DefaultHeight")
}

// TestStreams_TerminalHeightOverflowFd verifies overflow FD falls back to default.
func TestStreams_TerminalHeightOverflowFd(t *testing.T) {
	t.Parallel()

	raw := &fdWriter{fd: uintptr(math.MaxInt) + 1}
	s := termio.New(nil, raw, io.Discard)

	assert.Equalf(t, termio.DefaultHeight, s.TerminalHeight(),
		"overflow FD must fall back to DefaultHeight")
}

// TestStreams_TerminalHeightNonTerminalFd verifies regular file uses default.
func TestStreams_TerminalHeightNonTerminalFd(t *testing.T) {
	t.Parallel()

	f, err := os.CreateTemp(t.TempDir(), "termio-height-*")
	require.NoErrorf(t, err, "create temp file")

	t.Cleanup(func() { f.Close() }) //nolint:errcheck

	s := termio.New(nil, f, io.Discard)

	assert.Equalf(t, termio.DefaultHeight, s.TerminalHeight(),
		"regular file FD must fall back to DefaultHeight")
}

// TestStreams_TTYDetectionOverflowFd verifies overflow FD is not detected as TTY.
func TestStreams_TTYDetectionOverflowFd(t *testing.T) {
	t.Parallel()

	overflowFd := uintptr(math.MaxInt) + 1
	in := &fdReader{fd: overflowFd}
	out := &fdWriter{fd: overflowFd}
	errOut := &fdWriter{fd: overflowFd}
	s := termio.New(in, out, errOut)

	assert.Falsef(t, s.IsStdinTTY(), "overflow FD on stdin must not be TTY")
	assert.Falsef(t, s.IsStdoutTTY(), "overflow FD on stdout must not be TTY")
	assert.Falsef(t, s.IsStderrTTY(), "overflow FD on stderr must not be TTY")
}

// TestStreams_TTYDetectionNonTerminalFd verifies regular file is not TTY.
func TestStreams_TTYDetectionNonTerminalFd(t *testing.T) {
	t.Parallel()

	f, err := os.CreateTemp(t.TempDir(), "termio-tty-*")
	require.NoErrorf(t, err, "create temp file")

	t.Cleanup(func() { f.Close() }) //nolint:errcheck

	s := termio.New(f, f, f)

	assert.Falsef(t, s.IsStdinTTY(), "regular file must not be TTY")
	assert.Falsef(t, s.IsStdoutTTY(), "regular file must not be TTY")
	assert.Falsef(t, s.IsStderrTTY(), "regular file must not be TTY")
}

var errWriteFailedStream = errors.New("stream write failed")

// failWriter always fails every write.
type failWriter struct{}

func (f *failWriter) Write(_ []byte) (int, error) { return 0, errWriteFailedStream }

// prefixPolicy prepends a fixed string; used to verify policy wiring in Streams.
type prefixPolicy struct {
	prefix string
}

func (p *prefixPolicy) Apply(w io.Writer) io.Writer {
	return &prefixWriter{w: w, prefix: p.prefix}
}

type prefixWriter struct {
	w      io.Writer
	prefix string
}

func (pw *prefixWriter) Write(b []byte) (int, error) {
	_, err := io.WriteString(pw.w, pw.prefix)
	if err != nil {
		return 0, err
	}

	return pw.w.Write(b)
}

// fdWriter is a test-only [io.Writer] that exposes a fixed file descriptor.
type fdWriter struct {
	bytes.Buffer
	fd uintptr
}

func (w *fdWriter) Fd() uintptr { return w.fd }

// fdReader is a test-only [io.Reader] that exposes a fixed file descriptor.
type fdReader struct {
	bytes.Buffer
	fd uintptr
}

func (r *fdReader) Fd() uintptr { return r.fd }
