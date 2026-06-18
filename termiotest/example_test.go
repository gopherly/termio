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

	"gopherly.dev/termio/termiotest"
)

// ExampleNew returns buffer-backed streams and the underlying buffers for
// assertion. Writes to s.Out appear in the returned out buffer.
func ExampleNew() {
	s, _, out, _ := termiotest.New()

	fmt.Fprintln(s.Out, "hello") //nolint:errcheck
	fmt.Print(out.String())
	// Output:
	// hello
}

// ExampleNewTTY is like New but marks all three streams as TTYs so
// IsInteractive returns true.
func ExampleNewTTY() {
	s, in, out, errOut := termiotest.NewTTY()
	_ = in
	_ = out
	_ = errOut

	fmt.Println(s.IsInteractive())
	// Output:
	// true
}
