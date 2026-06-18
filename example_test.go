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
	"fmt"
	"io"

	"gopherly.dev/termio"
	"gopherly.dev/termio/colorprofile"
	"gopherly.dev/termio/termiotest"
)

// ExampleSystem demonstrates the zero-config constructor. In production
// code this connects to [os.Stdin], [os.Stdout], and [os.Stderr] with TTY
// detection performed against the real file descriptors.
func ExampleSystem() {
	// For reproducible example output we use termiotest.New() here; real
	// programs would call termio.System().
	s, _, out, _ := termiotest.New()

	fmt.Fprintln(s.Out, "hello from termio") //nolint:errcheck

	fmt.Print(out.String())
	// Output:
	// hello from termio
}

// ExampleNew shows how to supply your own readers and writers. This is the
// pattern used in tests and in code that needs to redirect output.
func ExampleNew() {
	s, _, out, errOut := termiotest.New()

	fmt.Fprintln(s.Out, "product output")        //nolint:errcheck
	fmt.Fprintln(s.ErrOut, "diagnostic message") //nolint:errcheck

	fmt.Print("out: ", out.String())
	fmt.Print("err: ", errOut.String())
	// Output:
	// out: product output
	// err: diagnostic message
}

// ExampleWriter_Err illustrates the per-stream sticky error. After a write
// fails, the error is latched on that Writer only. The other stream
// continues to work normally.
func ExampleWriter_Err() {
	errWriter := &alwaysFailWriter{}
	s := termio.New(nil, errWriter, io.Discard)

	fmt.Fprintln(s.Out, "this will fail") //nolint:errcheck

	fmt.Println("Out.Err:", s.Out.Err())
	fmt.Println("ErrOut.Err:", s.ErrOut.Err())
	// Output:
	// Out.Err: always fails
	// ErrOut.Err: <nil>
}

// ExampleWithColorPolicy shows how to inject a ColorPolicy at construction
// time. ANSI escape sequences are rewritten before they reach the output.
func ExampleWithColorPolicy() {
	buf := &bytes.Buffer{}
	env := []string{"NO_COLOR=1"}
	policy := colorprofile.Detect(buf, env)
	s := termio.New(nil, buf, io.Discard, termio.WithColorPolicy(policy))

	fmt.Fprintln(s.Out, "\x1b[32mhello\x1b[0m") //nolint:errcheck
	fmt.Print(buf.String())
	// Output:
	// hello
}

// ExampleStreams_TerminalWidth shows the fallback width for non-terminal
// streams such as buffer-backed writers used in tests.
func ExampleStreams_TerminalWidth() {
	s, in, out, errOut := termiotest.New()
	_ = in
	_ = out
	_ = errOut

	fmt.Println(s.TerminalWidth())
	// Output:
	// 80
}

// ExampleStreams_IsInteractive reports whether both stdin and stdout are
// terminals. termiotest.New returns non-TTY streams; NewTTY simulates a
// terminal session.
func ExampleStreams_IsInteractive() {
	plain, in, out, errOut := termiotest.New()
	_ = in
	_ = out
	_ = errOut

	tty, in2, out2, err2 := termiotest.NewTTY()
	_ = in2
	_ = out2
	_ = err2

	fmt.Println(plain.IsInteractive())
	fmt.Println(tty.IsInteractive())
	// Output:
	// false
	// true
}

// ExampleStreams_Err returns the first write error from Out or ErrOut. When
// both streams fail, the errors are joined.
func ExampleStreams_Err() {
	errWriter := &alwaysFailWriter{}
	s := termio.New(nil, errWriter, errWriter)

	fmt.Fprintln(s.Out, "fail")    //nolint:errcheck
	fmt.Fprintln(s.ErrOut, "fail") //nolint:errcheck

	fmt.Println(s.Err() != nil)
	// Output:
	// true
}

// alwaysFailWriter returns an error on every write call.
type alwaysFailWriter struct{}

func (a *alwaysFailWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("always fails") //nolint:err113
}
