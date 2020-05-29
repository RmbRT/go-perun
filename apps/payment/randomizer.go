// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by the Apache 2.0 license that can be found
// in the LICENSE file.

package payment

import (
	"math/rand"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/channel/test"
)

// Randomizer implements channel.test.AppRandomizer.
type Randomizer struct{}

var _ test.AppRandomizer = (*Randomizer)(nil)

// NewRandomApp always returns a payment app with the same address. Currently,
// one payment address has to be set at program startup.
func (*Randomizer) NewRandomApp(*rand.Rand) channel.App {
	return &App{AppDef()}
}

// NewRandomData returns NoData because a PaymentApp does not have data.
func (*Randomizer) NewRandomData(*rand.Rand) channel.Data {
	return new(NoData)
}
