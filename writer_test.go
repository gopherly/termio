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
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errWriteFailed = errors.New("write failed")

// failWriter always returns errWriteFailed on the first write.
type failWriter struct {
	bytes.Buffer
}

func (f *failWriter) Write(_ []byte) (int, error) { return 0, errWriteFailed }

// TestWriter_WriteForwardsBytes verifies bytes reach the underlying buffer.
func TestWriter_WriteForwardsBytes(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	w := newWriter(buf, nil)
	_, err := io.WriteString(w, "hello")

	require.NoErrorf(t, err, "write must not fail on healthy buffer")
	assert.Equalf(t, "hello", buf.String(), "bytes must reach underlying buffer")
}

// TestWriter_StickyErrorLatchesOnFirstFailure verifies the error is latched.
func TestWriter_StickyErrorLatchesOnFirstFailure(t *testing.T) {
	t.Parallel()

	w := newWriter(&failWriter{}, nil)

	_, err := io.WriteString(w, "trigger")

	require.ErrorIsf(t, err, errWriteFailed, "first write must return underlying error")
	assert.ErrorIsf(t, w.Err(), errWriteFailed, "Err() must return latched error")
}

// TestWriter_StickyErrorBlocksSubsequentWrites verifies later writes are blocked.
func TestWriter_StickyErrorBlocksSubsequentWrites(t *testing.T) {
	t.Parallel()

	w := newWriter(&failWriter{}, nil)
	io.WriteString(w, "trigger") //nolint:errcheck

	_, err := io.WriteString(w, "second")

	require.ErrorIsf(t, err, errWriteFailed, "subsequent write must return latched error")
}

// TestWriter_NoErrorInitially verifies Err() is nil before any failure.
func TestWriter_NoErrorInitially(t *testing.T) {
	t.Parallel()

	w := newWriter(&bytes.Buffer{}, nil)

	assert.NoErrorf(t, w.Err(), "Err() must be nil before any failed write")
}

// TestWriter_FdBuffer verifies that a plain buffer yields InvalidFd.
func TestWriter_FdBuffer(t *testing.T) {
	t.Parallel()

	w := newWriter(&bytes.Buffer{}, nil)

	assert.Equalf(t, InvalidFd, w.Fd(), "*bytes.Buffer must yield InvalidFd")
}

// TestWriter_FdPreservedThroughPolicy verifies Fd is captured before wrapping.
func TestWriter_FdPreservedThroughPolicy(t *testing.T) {
	t.Parallel()

	const want = uintptr(42)
	raw := &stubFdWriter{fd: want}

	w := newWriter(raw, noopPolicy{})

	assert.Equalf(t, want, w.Fd(), "Fd must be captured from raw before policy wraps it")
}

// TestWriter_PolicyAdapts verifies a ColorPolicy transforms written bytes.
func TestWriter_PolicyAdapts(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	prefix := "PREFIX:"

	p := &prefixPolicy{prefix: prefix}
	w := newWriter(buf, p)

	_, err := io.WriteString(w, "hello")

	require.NoErrorf(t, err, "policy write must not fail")
	assert.Equalf(t, prefix+"hello", buf.String(), "policy must transform the bytes")
}

// prefixPolicy is a test-only ColorPolicy that prepends a fixed string.
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
