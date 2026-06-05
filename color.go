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

// ColorPolicy adapts ANSI escape sequences for a specific terminal
// capability. Implementations intercept writes and strip, translate, or
// pass through color codes based on what the destination terminal supports.
//
// The core termio package does not ship a built-in implementation; callers
// who need color adaptation import [gopherly.dev/termio/colorprofile] and
// pass the result of [colorprofile.Detect] via [WithColorPolicy]. When no
// policy is configured, writes reach the raw stream unmodified.
//
// Implementors must be safe for concurrent use from multiple goroutines.
type ColorPolicy interface {
	// Apply wraps w with a writer that enforces this policy's color
	// translation rules. The returned writer must forward all writes to
	// w, potentially transforming ANSI sequences in transit. Apply must
	// not retain w beyond the lifetime of the returned writer.
	Apply(w io.Writer) io.Writer
}
