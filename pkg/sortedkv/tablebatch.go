// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by the Apache 2.0 license that can be found
// in the LICENSE file.

package sortedkv

// tableBatch is a wrapper around a Database Batch with a key prefix. All
// Writer operations are automatically prefixed.
type tableBatch struct {
	Batch
	prefix string
}

func (b *tableBatch) pkey(key string) string {
	return b.prefix + key
}

// Put puts a value into a table batch.
func (b *tableBatch) Put(key, value string) error {
	return b.Batch.Put(b.pkey(key), value)
}

// Put puts a value into a table batch.
func (b *tableBatch) PutBytes(key string, value []byte) error {
	return b.Batch.PutBytes(b.pkey(key), value)
}

// Delete deletes a value from a table batch.
func (b *tableBatch) Delete(key string) error {
	return b.Batch.Delete(b.pkey(key))
}
