// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

// Package rand implements a cryptographically secure
// pseudorandom number generator.
package rand

import (
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/binary"
	"io"
	"sync"

	"github.com/tmthrgd/chacha20"
)

const (
	// SeedSize is the required length of seed passed
	// to New.
	SeedSize = chacha20.KeySize

	// budget is the number of bytes that can be
	// generated from a freshly seeded generator
	// before a reseed.
	//
	// The Draft cipher has a 64-bit counter so can
	// handle 2^70-1 before rolling over. Inspite of
	// this, a more conservative 2^30 is used.
	//
	// On top of the budget a further SeedSize bytes
	// will be used to reseed the generator.
	budget = 1 << 30
)

var (
	// Reader is a global, shared instance of a
	// cryptographically strong pseudo-random
	// generator.
	//
	// The seed is read from crypto/rand.Reader.
	Reader io.Reader = &seedOnFirstRead{}

	// zero is an array of zeros (length SHOULD be
	// mutltiple of 128) that will be XORd to give
	// the keystream directly. (K XOR 0 = K).
	//
	// The length of zero sets the maximum number of
	// bytes that can be read in a single call to
	// XORKeyStream (along with (*reader).budget).
	zero [1024 * 1024]byte
)

// Read is a helper function that calls Reader.Read using
// io.ReadFull. On return, n == len(b) if and only if
// err == nil.
func Read(b []byte) (n int, err error) {
	return io.ReadFull(Reader, b)
}

// New returns a new pseudorandom generator with the given
// seed. If seed == nil, the generator seeds itself by
// reading from crypto/rand.Reader. seed must be SeedSize
// bytes long.
//
// The Read method on the returned reader always returns
// the full amount asked for, or else it returns an error.
//
// The generator uses ChaCha20 reseeding after every
// 1 GB of generated data.
//
// The generator is deterministic for a given seed.
func New(seed []byte) (io.Reader, error) {
	if seed == nil {
		seed = make([]byte, SeedSize)

		if _, err := crand.Read(seed); err != nil {
			return nil, err
		}
	}

	var nonce [chacha20.DraftNonceSize]byte

	c, err := chacha20.NewDraft(seed, nonce[:])
	if err != nil {
		return nil, err
	}

	return &reader{
		cipher: c,
		budget: budget,
	}, nil
}

type reader struct {
	mu sync.Mutex

	cipher  cipher.Stream
	budget  int
	counter uint64
}

func (r *reader) Read(b []byte) (n int, err error) {
	n = len(b)

	r.mu.Lock()

	for len(b) != 0 {
		if r.budget == 0 {
			var key [chacha20.KeySize]byte
			r.cipher.XORKeyStream(key[:], key[:])

			var nonce [chacha20.DraftNonceSize]byte
			binary.LittleEndian.PutUint64(nonce[:], r.counter+1)

			var c cipher.Stream
			c, err = chacha20.NewDraft(key[:], nonce[:])
			if err != nil {
				break
			}

			r.cipher = c
			r.budget = budget
			r.counter++
		}

		todo := len(b)
		if todo > len(zero) {
			todo = len(zero)
		}
		if todo > r.budget {
			todo = r.budget
		}

		r.cipher.XORKeyStream(b[:todo], zero[:todo])

		r.budget -= todo
		b = b[todo:]
	}

	r.mu.Unlock()

	n -= len(b)
	return
}

// This (partially) works around
// https://github.com/golang/go/issues/11833 by only
// seeding the generator upon the first call to Read.
type seedOnFirstRead struct {
	once sync.Once
	r    io.Reader
	err  error
}

func (r *seedOnFirstRead) Read(b []byte) (n int, err error) {
	r.once.Do(r.seed)
	if r.err != nil {
		return 0, r.err
	}

	return r.r.Read(b)
}

func (r *seedOnFirstRead) seed() {
	r.r, r.err = New(nil)
}
