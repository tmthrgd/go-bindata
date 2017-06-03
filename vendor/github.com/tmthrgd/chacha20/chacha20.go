// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

// Package chacha20 provides an AVX/AVX2/pure-Go implementation of ChaCha20, a
// fast, secure stream cipher.
//
// From Bernstein, Daniel J. "ChaCha, a variant of Salsa20." Workshop Record of
// SASC. 2008. (http://cr.yp.to/chacha/chacha-20080128.pdf):
//
//	ChaCha8 is a 256-bit stream cipher based on the 8-round cipher Salsa20/8.
//	The changes from Salsa20/8 to ChaCha8 are designed to improve diffusion per
//	round, conjecturally increasing resistance to cryptanalysis, while
//	preserving -- and often improving -- time per round. ChaCha12 and ChaCha20
//	are analogous modiﬁcations of the 12-round and 20-round ciphers Salsa20/12
//	and Salsa20/20. This paper presents the ChaCha family and explains the
//	differences between Salsa20 and ChaCha.
//
// For more information, see http://cr.yp.to/chacha.html
package chacha20

import (
	"crypto/cipher"
	"errors"
)

const (
	// KeySize is the length of ChaCha20 keys, in bytes.
	KeySize = 32

	// NonceSize is the length of ChaCha20 nonces, in bytes.
	//
	// In most cases either RFCNonceSize or DraftNonceSize should
	// be used instead.
	//
	// This is maintained for compatibility reasons.
	NonceSize = DraftNonceSize

	// RFCNonceSize is the length of ChaCha20-RFC nonces, in bytes.
	RFCNonceSize = 12

	// DraftNonceSize is the length of ChaCha20-draft nonces, in bytes.
	DraftNonceSize = 8

	// XNonceSize is the length of XChaCha20 nonces, in bytes.
	XNonceSize = 24
)

var (
	// ErrInvalidKey is returned when the provided key is not KeySize bytes long.
	ErrInvalidKey = errors.New("invalid key length")

	// ErrInvalidNonce is returned when the provided nonce is not RFCNonceSize,
	// DraftNonceSize or XNonceSize bytes long.
	ErrInvalidNonce = errors.New("invalid nonce length")
)

// New creates and returns a new cipher.Stream. The key argument must be 256
// bits long, and the nonce argument must be either 64, 96 or 192 bits long.
// The nonce must be randomly generated or used only once. If the nonce
// argument is 64 bits long, New behaves like NewDraft. If the nonce argument
// is 96 bits long, New behaves like NewRFC. If the nonce argument is 192 bits
// long New behaves like NewXChaCha.
//
// In most cases either NewRFC, NewDraft or NewXChaCha should be used instead.
func New(key, nonce []byte) (cipher.Stream, error) {
	switch len(nonce) {
	case XNonceSize:
		return NewXChaCha(key, nonce)
	case RFCNonceSize:
		return NewRFC(key, nonce)
	default:
		return NewDraft(key, nonce)
	}
}
