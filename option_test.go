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
	"testing"

	"github.com/stretchr/testify/assert"
)

// noopPolicy is a minimal ColorPolicy implementation for option tests.
type noopPolicy struct{}

func (noopPolicy) Apply(w io.Writer) io.Writer { return w }

// TestWithColorPolicy_SetsPolicy verifies the option stores a policy.
func TestWithColorPolicy_SetsPolicy(t *testing.T) {
	t.Parallel()

	p := noopPolicy{}
	cfg := &config{}
	WithColorPolicy(p)(cfg)

	assert.Equalf(t, p, cfg.colorPolicy, "WithColorPolicy must store the supplied policy")
}

// TestWithColorPolicy_Nil verifies that a nil policy clears the config.
func TestWithColorPolicy_Nil(t *testing.T) {
	t.Parallel()

	cfg := &config{colorPolicy: noopPolicy{}}
	WithColorPolicy(nil)(cfg)

	assert.Nilf(t, cfg.colorPolicy, "WithColorPolicy(nil) must clear the policy")
}

// TestApplyAll_Empty verifies that no options leave config zero-valued.
func TestApplyAll_Empty(t *testing.T) {
	t.Parallel()

	cfg := &config{}
	applyAll(cfg, nil)

	assert.Nilf(t, cfg.colorPolicy, "applyAll with no opts must leave config zero-valued")
}

// TestApplyAll_LastWins verifies that the last option takes precedence.
func TestApplyAll_LastWins(t *testing.T) {
	t.Parallel()

	first := noopPolicy{}
	second := noopPolicy{}

	cfg := &config{}
	applyAll(cfg, []Option{
		WithColorPolicy(first),
		WithColorPolicy(second),
	})

	assert.Equalf(t, second, cfg.colorPolicy, "last option must win")
}
