// Copyright (c) 2020 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by the Apache 2.0 license that can be found
// in the LICENSE file.

package test

import (
	"context"
	"math/rand"

	"github.com/stretchr/testify/require"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/channel/persistence"
	ctest "perun.network/go-perun/channel/test"
	"perun.network/go-perun/wallet"
	wtest "perun.network/go-perun/wallet/test"
	"perun.network/go-perun/wire"
)

// Channel is a wrapper around a persisted channel and its participants, as well
// as the associated persister and restorer.
type Channel struct {
	accounts []wallet.Account
	peers    []wire.Address
	*persistence.StateMachine

	pr  persistence.PersistRestorer
	ctx context.Context
}

// NewRandomChannel creates a random channel with the requested persister and
// restorer, as well as the selected peer addresses for the participants (other
// than the owner's). The owner's index in the channel participants can be
// controlled via the 'user' argument. The wallet accounts and addresses used by
// the participants are generated randomly.
// The persister is notified and called to persist the new channel before it is
// returned.
func NewRandomChannel(
	ctx context.Context,
	t require.TestingT,
	pr persistence.PersistRestorer,
	user channel.Index,
	peers []wire.Address,
	rng *rand.Rand) (c *Channel) {

	accs, parts := wtest.NewRandomAccounts(rng, len(peers))
	params := ctest.NewRandomParams(rng, ctest.WithParts(parts...))
	csm, err := channel.NewStateMachine(accs[0], *params)
	require.NoError(t, err)

	sm := persistence.FromStateMachine(csm, pr)
	c = &Channel{
		accounts:     accs,
		peers:        peers,
		StateMachine: &sm,
		pr:           pr,
		ctx:          ctx,
	}

	require.NoError(t, pr.ChannelCreated(ctx, c.StateMachine, c.peers))

	return
}

// AssertPersisted reads the channel state from the restorer and compares it
// to the actual channel state. If an error occurs while restoring the channel
// or if the restored channel does not match the actual channel state, then the
// test fails.
func (c *Channel) AssertPersisted(ctx context.Context, t require.TestingT) {
	ch, err := c.pr.RestoreChannel(ctx, c.ID())
	require.NoError(t, err)
	require.NotNil(t, ch)
	c.RequireEqual(t, ch)
}

// RequireEqual asserts that the channel is equal to the provided channel state.
func (c *Channel) RequireEqual(t require.TestingT, ch channel.Source) {
	require.Equal(t, c.Idx(), ch.Idx(), "Idx")
	require.Equal(t, c.Params(), ch.Params(), "Params")
	requireEqualStagingTX(t, c.StagingTX(), ch.StagingTX())
	require.Equal(t, c.CurrentTX(), ch.CurrentTX(), "CurrentTX")
	require.Equal(t, c.Phase(), ch.Phase(), "Phase")
}

// EqualStagingLoose is a test for loose equality between two staging states,
// where it is allowed for signatures to be a nil slice iff the transaction
// which it is compared to also has a nil slice OR a slice of nil sigs.
func requireEqualStagingTX(t require.TestingT, expected, actual channel.Transaction) {
	require.Equal(t, expected.State, actual.State, "StagingTX.State")
	requireEqualSigs(t, expected.Sigs, actual.Sigs)
}

func requireEqualSigs(t require.TestingT, expected, actual []wallet.Sig) {
	if expected == nil && actual == nil {
		return
	}
	actualNil := isNilSigs(actual)
	expectedNil := isNilSigs(expected)
	if (expected == nil && actualNil) ||
		(expectedNil && actual == nil) {
		return
	}
	if actualNil && expectedNil {
		if len(expected) != len(actual) {
			t.FailNow()
		}
		return
	}
	require.Equal(t, expected, actual, "StagingTX.Sigs")
}

func isNilSigs(s []wallet.Sig) bool {
	for _, el := range s {
		if el != nil {
			return false
		}
	}
	return true
}

// Init calls Init on the state machine and then checks the persistence.
func (c *Channel) Init(t require.TestingT, rng *rand.Rand) {
	initAlloc := *ctest.NewRandomAllocation(rng, ctest.WithNumParts(len(c.accounts)))
	initData := channel.NewMockOp(channel.OpValid)
	err := c.StateMachine.Init(nil, initAlloc, initData)
	require.NoError(t, err)
	c.AssertPersisted(c.ctx, t)
}

// EnableInit calls EnableInit on the state machine and then checks the persistence.
func (c *Channel) EnableInit(t require.TestingT) {
	err := c.StateMachine.EnableInit(c.ctx)
	require.NoError(t, err)
	c.AssertPersisted(c.ctx, t)
}

// SignAll signs the current staged state by all parties.
func (c *Channel) SignAll(t require.TestingT) {
	_, err := c.Sig(nil) // trigger local signing
	require.NoError(t, err)
	c.AssertPersisted(c.ctx, t)
	// remote signers
	for i := range c.accounts {
		sig, err := channel.Sign(c.accounts[i], c.Params(), c.StagingState())
		require.NoError(t, err)
		c.AddSig(nil, channel.Index(i), sig)
		c.AssertPersisted(c.ctx, t)
	}
}

// SetFunded calls SetFunded on the state machine and then checks the persistence.
func (c *Channel) SetFunded(t require.TestingT) {
	require.NoError(t, c.StateMachine.SetFunded(c.ctx))
	c.AssertPersisted(c.ctx, t)
}

// Update calls Update on the state machine and then checks the persistence.
func (c *Channel) Update(t require.TestingT, state *channel.State, idx channel.Index) error {
	err := c.StateMachine.Update(c.ctx, state, idx)
	c.AssertPersisted(c.ctx, t)
	return err
}

// EnableUpdate calls EnableUpdate on the state machine and then checks the persistence.
func (c *Channel) EnableUpdate(t require.TestingT) {
	require.NoError(t, c.StateMachine.EnableUpdate(c.ctx))
	c.AssertPersisted(c.ctx, t)
}

// EnableFinal calls EnableFinal on the state machine and then checks the persistence.
func (c *Channel) EnableFinal(t require.TestingT) {
	require.NoError(t, c.StateMachine.EnableFinal(c.ctx))
	c.AssertPersisted(c.ctx, t)
}

// DiscardUpdate calls DiscardUpdate on the state machine and then checks the persistence.
func (c *Channel) DiscardUpdate(t require.TestingT) {
	require.NoError(t, c.StateMachine.DiscardUpdate(c.ctx))
	c.AssertPersisted(c.ctx, t)
}

// SetRegistering calls SetRegistering on the state machine and then checks the persistence.
func (c *Channel) SetRegistering(t require.TestingT) {
	require.NoError(t, c.StateMachine.SetRegistering(c.ctx))
	c.AssertPersisted(c.ctx, t)
}

// SetRegistered calls SetRegistered on the state machine and then checks the persistence.
func (c *Channel) SetRegistered(t require.TestingT, r *channel.RegisteredEvent) {
	require.NoError(t, c.StateMachine.SetRegistered(c.ctx, r))
	c.AssertPersisted(c.ctx, t)
}

// SetWithdrawing calls SetWithdrawing on the state machine and then checks the persistence.
func (c *Channel) SetWithdrawing(t require.TestingT) {
	require.NoError(t, c.StateMachine.SetWithdrawing(c.ctx))
	c.AssertPersisted(c.ctx, t)
}

// SetWithdrawn calls SetWithdrawn on the state machine and then checks the persistence.
func (c *Channel) SetWithdrawn(t require.TestingT) {
	require.NoError(t, c.StateMachine.SetWithdrawn(c.ctx))
	c.AssertPersisted(c.ctx, t)
}

// Settle removes the channels data from the db and checks whether it really was
// removed from said db or not.
func (c *Channel) Settle(t require.TestingT) {
	c.pr.ChannelRemoved(c.ctx, c.ID())
	rc, err := c.pr.RestoreChannel(c.ctx, c.ID())
	require.Error(t, err, "restoring of a non-existing channel")
	require.Nil(t, rc, "restoring of a non-existing channel")
}
