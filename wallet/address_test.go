// Copyright (c) 2020 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by the Apache 2.0 license that can be found
// in the LICENSE file.

package wallet_test

import (
	"io"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	_ "perun.network/go-perun/backend/sim/wallet"
	iotest "perun.network/go-perun/pkg/io/test"
	pkgtest "perun.network/go-perun/pkg/test"
	"perun.network/go-perun/wallet"
	wallettest "perun.network/go-perun/wallet/test"
)

type testAddresses struct {
	addrs wallet.AddressesWithLen
}

func (t *testAddresses) Encode(w io.Writer) error {
	return t.addrs.Encode(w)
}

func (t *testAddresses) Decode(r io.Reader) error {
	return t.addrs.Decode(r)
}

func TestAddresses_Serializer(t *testing.T) {
	rng := rand.New(rand.NewSource(0xC00FED))

	addrs := wallettest.NewRandomAddresses(rng, 0)
	iotest.GenericSerializerTest(t, &testAddresses{addrs})

	addrs = wallettest.NewRandomAddresses(rng, 1)
	iotest.GenericSerializerTest(t, &testAddresses{addrs})

	addrs = wallettest.NewRandomAddresses(rng, 5)
	iotest.GenericSerializerTest(t, &testAddresses{addrs})
}

func TestAddrKey_Equals(t *testing.T) {
	rng := pkgtest.Prng(t)
	addrs := wallettest.NewRandomAddresses(rng, 10)

	// Test all properties of an equivalence relation.
	for i, a := range addrs {
		for j, b := range addrs {
			// Symmetry.
			require.Equal(t, wallet.Key(a).Equals(b), wallet.Key(b).Equals(a))
			// Test that Equals is equivalent to ==.
			require.Equal(t, wallet.Key(a).Equals(b), wallet.Key(a) == wallet.Key(b))
			// Test that it is not trivially set to true or false.
			require.Equal(t, i == j, wallet.Key(a).Equals(b))
			// Transitivity.
			for _, c := range addrs {
				if wallet.Key(a).Equals(b) && wallet.Key(b).Equals(c) {
					require.True(t, wallet.Key(a).Equals(c))
				}
			}
		}
		// Reflexivity.
		require.True(t, wallet.Key(a).Equals(a))
	}
}

func TestAddrKey(t *testing.T) {
	rng := pkgtest.Prng(t)
	addrs := wallettest.NewRandomAddresses(rng, 10)

	for _, a := range addrs {
		// Test that Key and FromKey are dual to each other.
		require.Equal(t, wallet.Key(a), wallet.Key(wallet.FromKey(wallet.Key(a))))
		// Test that FromKey returns the correct Address.
		require.True(t, a.Equals(wallet.FromKey(wallet.Key(a))))
	}
}
