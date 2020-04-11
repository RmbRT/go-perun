// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by a MIT-style license that can be found in
// the LICENSE file.

package channel_test

import (
	"math/rand"
	"testing"

	_ "perun.network/go-perun/backend/sim" // backend init
	"perun.network/go-perun/channel/test"
	iotest "perun.network/go-perun/pkg/io/test"
)

func TestStateSerialization(t *testing.T) {
	rng := rand.New(rand.NewSource(1337))

	app := test.NewRandomApp(rng)
	params := test.NewRandomParams(rng, app.Def())
	state := test.NewRandomState(rng, params)

	iotest.GenericSerializerTest(t, state)
}
