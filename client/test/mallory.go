// Copyright (c) 2020 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by a MIT-style license that can be found in
// the LICENSE file.

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"perun.network/go-perun/client"
)

// Mallory is a test client role. She proposes the new channel.
type Mallory struct {
	Proposer
}

// NewMallory creates a new party that executes the Mallory protocol.
func NewMallory(setup RoleSetup, t *testing.T) *Mallory {
	return &Mallory{Proposer: *NewProposer(setup, t, 3)}
}

// Execute executes the Mallory protocol.
func (r *Mallory) Execute(cfg ExecConfig) {
	r.Proposer.Execute(cfg, r.exec)
}

func (r *Mallory) exec(cfg ExecConfig, ch *paymentChannel) {
	assert := assert.New(r.t)
	we, _ := r.Idxs(cfg.PeerAddrs)
	// AdjudicatorReq for version 0
	req0 := client.NewTestChannel(ch.Channel).AdjudicatorReq()

	// 1st stage - channel controller set up
	r.waitStage()

	// Mallory sends some updates to Carol
	for i := 0; i < cfg.NumUpdates[we]; i++ {
		ch.sendTransfer(cfg.TxAmounts[we], fmt.Sprintf("Mallory#%d", i))
	}
	// 2nd stage - txs sent
	r.waitStage()

	// Register version 0 AdjudicatorReq
	challengeDuration := time.Duration(ch.Channel.Params().ChallengeDuration) * time.Second
	regCtx, regCancel := context.WithTimeout(context.Background(), r.timeout)
	defer regCancel()
	r.log.Debug("Registering version 0 state.")
	reg0, err := r.setup.Adjudicator.Register(regCtx, req0)
	assert.NoError(err)
	assert.NotNil(reg0)
	r.log.Debugln("<Registered> ver 0: ", reg0)

	// within the challenge duration, Carol should refute.
	subCtx, subCancel := context.WithTimeout(context.Background(), r.timeout+challengeDuration)
	defer subCancel()
	sub, err := r.setup.Adjudicator.SubscribeRegistered(subCtx, ch.Params())

	// 3rd stage - wait until Carol has refuted
	r.waitStage()

	assert.True(reg0.Timeout.IsElapsed(subCtx),
		"Carol's refutation should already have progressed past the timeout.")
	reg := sub.Next() // should be event caused by Carol's refutation.
	assert.NoError(sub.Close())
	assert.NoError(sub.Err())
	assert.NotNil(reg)
	r.log.Debugln("<Registered> refuted: ", reg)
	if reg != nil {
		assert.Equal(ch.State().Version, reg.Version, "expected refutation with current version")
		waitCtx, waitCancel := context.WithTimeout(context.Background(), r.timeout+challengeDuration)
		defer waitCancel()
		// refutation increased the timeout.
		assert.NoError(reg.Timeout.Wait(waitCtx))
	}

	wdCtx, wdCancel := context.WithTimeout(context.Background(), r.timeout)
	defer wdCancel()
	err = r.setup.Adjudicator.Withdraw(wdCtx, req0)
	assert.Error(err, "withdrawing should fail because Carol should have refuted.")

	// settling current version should work
	ch.settleChan()
}
