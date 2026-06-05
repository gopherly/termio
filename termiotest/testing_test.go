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

package termiotest_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopherly.dev/termio/termiotest"
)

// TestNew_BuffersAreLinked verifies writes reach the returned buffers.
func TestNew_BuffersAreLinked(t *testing.T) {
	t.Parallel()

	s, _, out, errOut := termiotest.New()

	_, err := fmt.Fprintln(s.Out, "stdout-line")
	require.NoErrorf(t, err, "Fprintln to Out must succeed")

	_, err = fmt.Fprintln(s.ErrOut, "stderr-line")
	require.NoErrorf(t, err, "Fprintln to ErrOut must succeed")

	assert.Equalf(t, "stdout-line\n", out.String(), "Out writes must appear in out buffer")
	assert.Equalf(t, "stderr-line\n", errOut.String(), "ErrOut writes must appear in errOut buffer")
}

// TestNew_AllNonTTY verifies New returns non-TTY streams.
func TestNew_AllNonTTY(t *testing.T) {
	t.Parallel()

	s, _, _, _ := termiotest.New() //nolint:dogsled

	assert.Falsef(t, s.IsStdinTTY(), "New must return non-TTY stdin")
	assert.Falsef(t, s.IsStdoutTTY(), "New must return non-TTY stdout")
	assert.Falsef(t, s.IsStderrTTY(), "New must return non-TTY stderr")
	assert.Falsef(t, s.IsInteractive(), "New must not be interactive")
}

// TestNew_InBufferIsReadable verifies s.In reads from the in buffer.
func TestNew_InBufferIsReadable(t *testing.T) {
	t.Parallel()

	s, in, _, _ := termiotest.New()

	_, err := fmt.Fprint(in, "input-data")
	require.NoErrorf(t, err, "writing to in buffer must succeed")

	got, err := io.ReadAll(s.In)

	require.NoErrorf(t, err, "reading from s.In must succeed")
	assert.Equalf(t, "input-data", string(got), "s.In must read from in buffer")
}

// TestNewTTY_AllTTY verifies NewTTY returns TTY-flagged streams.
func TestNewTTY_AllTTY(t *testing.T) {
	t.Parallel()

	s, _, _, _ := termiotest.NewTTY() //nolint:dogsled

	assert.Truef(t, s.IsStdinTTY(), "NewTTY must return TTY stdin")
	assert.Truef(t, s.IsStdoutTTY(), "NewTTY must return TTY stdout")
	assert.Truef(t, s.IsStderrTTY(), "NewTTY must return TTY stderr")
	assert.Truef(t, s.IsInteractive(), "NewTTY must be interactive")
}

// TestNewTTY_BuffersStillWritable verifies TTY flag does not break writes.
func TestNewTTY_BuffersStillWritable(t *testing.T) {
	t.Parallel()

	s, _, out, _ := termiotest.NewTTY()

	_, err := fmt.Fprint(s.Out, "tty-write")

	require.NoErrorf(t, err, "write to TTY-flagged Out must succeed")
	assert.Equalf(t, "tty-write", out.String(), "TTY flag must not affect buffer writes")
}
