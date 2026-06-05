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

// Option configures a [Streams] at construction time via [New] or [System].
// Use the With* helpers to build options; do not implement Option directly.
type Option func(*config)

// config holds the resolved functional options for [New] and [System].
type config struct {
	colorPolicy ColorPolicy
}

// applyAll applies all opts to cfg in order.
func applyAll(cfg *config, opts []Option) {
	for _, o := range opts {
		o(cfg)
	}
}

// WithColorPolicy sets the [ColorPolicy] applied to Out and ErrOut during
// [Writer] construction. When not set (or set to nil), writes reach the raw
// stream unmodified.
//
// Example using the bundled colorprofile adapter:
//
//	import "gopherly.dev/termio/colorprofile"
//
//	s := termio.System(
//	    termio.WithColorPolicy(colorprofile.Detect(os.Stdout, os.Environ())),
//	)
func WithColorPolicy(p ColorPolicy) Option {
	return func(c *config) {
		c.colorPolicy = p
	}
}
