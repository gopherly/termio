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
	"io"
	"os"
)

// InvalidFd is the sentinel value returned by [Writer.Fd] when the
// underlying stream is not backed by a real file descriptor (for example,
// when a [bytes.Buffer] is supplied in tests). It equals ^uintptr(0), the
// maximum uintptr value, which never refers to a valid OS descriptor on any
// supported platform.
const InvalidFd = ^uintptr(0)

// fileWriter is the interface any TTY-aware library should accept: an
// [io.Writer] that also exposes its file descriptor. Preserving Fd()
// through the write chain lets downstream libraries (bubbletea, glamour,
// lipgloss) detect the terminal even when they receive a wrapped writer.
type fileWriter interface {
	io.Writer
	Fd() uintptr
}

// fdOf reports the file descriptor of v when v is an *[os.File] or already
// implements interface{ Fd() uintptr }; otherwise it returns [InvalidFd].
func fdOf(v any) uintptr {
	switch f := v.(type) {
	case *os.File:
		if f == nil {
			return InvalidFd
		}

		return f.Fd()
	case interface{ Fd() uintptr }:
		return f.Fd()
	default:
		return InvalidFd
	}
}
