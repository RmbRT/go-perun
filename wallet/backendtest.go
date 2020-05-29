// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by the Apache 2.0 license that can be found
// in the LICENSE file.

package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SetBackendTest is a generic test to test that the wallet backend is set correctly.
func SetBackendTest(t *testing.T) {
	assert.Panics(t, func() { SetBackend(nil) }, "nil backend set should panic")
	require.NotNil(t, backend, "backend should be already set by init()")
	assert.Panics(t, func() { SetBackend(backend) }, "setting a backend twice should panic")
}
