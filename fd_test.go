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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFdOf_NilFile verifies that a nil *[os.File] returns InvalidFd.
func TestFdOf_NilFile(t *testing.T) {
	t.Parallel()

	assert.Equalf(t, InvalidFd, fdOf((*os.File)(nil)), "nil *os.File must return InvalidFd")
}

// TestFdOf_Buffer verifies that a *[bytes.Buffer] returns InvalidFd.
func TestFdOf_Buffer(t *testing.T) {
	t.Parallel()

	assert.Equalf(t, InvalidFd, fdOf(&bytes.Buffer{}), "*bytes.Buffer must return InvalidFd")
}

// TestFdOf_FdIface verifies that a type implementing Fd() is forwarded.
func TestFdOf_FdIface(t *testing.T) {
	t.Parallel()

	const want = uintptr(42)
	w := &stubFdWriter{fd: want}

	assert.Equalf(t, want, fdOf(w), "Fd() implementation must be forwarded")
}

// TestFdOf_RealFile verifies that a real *[os.File] returns a valid FD.
func TestFdOf_RealFile(t *testing.T) {
	t.Parallel()

	f, err := os.CreateTemp(t.TempDir(), "termio-fd-test-*")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}

	t.Cleanup(func() { f.Close() }) //nolint:errcheck

	got := fdOf(f)

	assert.NotEqualf(t, InvalidFd, got, "real *os.File must not return InvalidFd")
}

// stubFdWriter is a minimal fileWriter used only in tests.
type stubFdWriter struct {
	bytes.Buffer
	fd uintptr
}

func (s *stubFdWriter) Fd() uintptr { return s.fd }
