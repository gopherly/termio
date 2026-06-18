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

//go:build !windows

package termio_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"

	"gopherly.dev/termio"
)

func openPTY(t *testing.T) *os.File {
	t.Helper()

	ptmx, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	require.NoErrorf(t, err, "open /dev/ptmx")

	t.Cleanup(func() { ptmx.Close() }) //nolint:errcheck

	return ptmx
}

func setWinsize(t *testing.T, f *os.File, cols, rows uint16) {
	t.Helper()

	ws := unix.Winsize{Col: cols, Row: rows}
	err := unix.IoctlSetWinsize(int(f.Fd()), unix.TIOCSWINSZ, &ws)

	require.NoErrorf(t, err, "TIOCSWINSZ must succeed")
}

// TestStreams_TTYDetectionRealTerminal verifies TTY detection against a PTY.
func TestStreams_TTYDetectionRealTerminal(t *testing.T) {
	t.Parallel()

	pty := openPTY(t)
	s := termio.New(pty, pty, pty)

	assert.Truef(t, s.IsStdinTTY(), "PTY must be detected as TTY for stdin")
	assert.Truef(t, s.IsStdoutTTY(), "PTY must be detected as TTY for stdout")
	assert.Truef(t, s.IsStderrTTY(), "PTY must be detected as TTY for stderr")
	assert.Truef(t, s.IsInteractive(), "PTY-backed streams must be interactive")
}

// TestStreams_TerminalWidthRealTerminal verifies width detection against a PTY.
func TestStreams_TerminalWidthRealTerminal(t *testing.T) {
	t.Parallel()

	pty := openPTY(t)
	setWinsize(t, pty, 120, 40)

	s := termio.New(nil, pty, io.Discard)

	assert.Equalf(t, 120, s.TerminalWidth(),
		"PTY-backed TerminalWidth must reflect the configured size")
}

// TestStreams_TerminalHeightRealTerminal verifies height detection against a PTY.
func TestStreams_TerminalHeightRealTerminal(t *testing.T) {
	t.Parallel()

	pty := openPTY(t)
	setWinsize(t, pty, 120, 40)

	s := termio.New(nil, pty, io.Discard)

	assert.Equalf(t, 40, s.TerminalHeight(),
		"PTY-backed TerminalHeight must reflect the configured size")
}
