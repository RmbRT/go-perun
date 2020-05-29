// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by the Apache 2.0 license that can be found
// in the LICENSE file.

package client_test

import (
	"math/rand"

	"github.com/sirupsen/logrus"

	"perun.network/go-perun/apps/payment"
	plogrus "perun.network/go-perun/log/logrus"
	wallettest "perun.network/go-perun/wallet/test"
)

func init() {
	plogrus.Set(logrus.WarnLevel, &logrus.TextFormatter{ForceColors: true})

	// Eth client tests use the payment app for now...
	// TODO: This has to be set to the deployed app contract (or counterfactual
	// address of it) when we start using it in tests.
	// Use random seed that should be different from other seeds used in tests.
	rng := rand.New(rand.NewSource(0x280a0f350eec))
	appDef := wallettest.NewRandomAddress(rng)
	payment.SetAppDef(appDef) // payment app address has to be set once at startup
}
