// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by the Apache 2.0 license that can be found
// in the LICENSE file.

package wire

import (
	"context"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	pkgtest "perun.network/go-perun/pkg/test"
	wallettest "perun.network/go-perun/wallet/test"
)

func TestAuthResponseMsg(t *testing.T) {
	rng := rand.New(rand.NewSource(1337))
	TestMsg(t, NewAuthResponseMsg(wallettest.NewRandomAccount(rng)))
}

func TestExchangeAddrs_ConnFail(t *testing.T) {
	rng := rand.New(rand.NewSource(0xDDDDDEDE))
	a, _ := newPipeConnPair()
	a.Close()
	addr, err := ExchangeAddrs(context.Background(), wallettest.NewRandomAccount(rng), a)
	assert.Nil(t, addr)
	assert.Error(t, err)
}

func TestExchangeAddrs_Success(t *testing.T) {
	rng := rand.New(rand.NewSource(0xfedd))
	conn0, conn1 := newPipeConnPair()
	defer conn0.Close()
	account0, account1 := wallettest.NewRandomAccount(rng), wallettest.NewRandomAccount(rng)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer conn1.Close()

		recvAddr0, err := ExchangeAddrs(context.Background(), account1, conn1)
		assert.NoError(t, err)
		assert.True(t, recvAddr0.Equals(account0.Address()))
	}()

	recvAddr1, err := ExchangeAddrs(context.Background(), account0, conn0)
	assert.NoError(t, err)
	assert.True(t, recvAddr1.Equals(account1.Address()))

	wg.Wait()
}

func TestExchangeAddrs_Timeout(t *testing.T) {
	rng := rand.New(rand.NewSource(0xDDDDDeDe))
	a, _ := newPipeConnPair()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	pkgtest.AssertTerminates(t, 2*timeout, func() {
		addr, err := ExchangeAddrs(ctx, wallettest.NewRandomAccount(rng), a)
		assert.Nil(t, addr)
		assert.Error(t, err)
	})
}

func TestExchangeAddrs_BogusMsg(t *testing.T) {
	rng := rand.New(rand.NewSource(0xcafe))
	acc := wallettest.NewRandomAccount(rng)
	conn := newMockConn(nil)
	conn.recvQueue <- NewPingMsg()
	addr, err := ExchangeAddrs(context.Background(), acc, conn)

	assert.Error(t, err, "ExchangeAddrs should error when peer sends a non-AuthResponseMsg")
	assert.Nil(t, addr)
}
