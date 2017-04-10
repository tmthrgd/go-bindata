// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

// +build !amd64 gccgo appengine

package chacha20

import (
	"crypto/cipher"

	"github.com/tmthrgd/chacha20/internal/ref"
)

const useRef = true

var useAVX, useAVX2 = false, false

// NewRFC creates and returns a new cipher.Stream. The key argument must be 256
// bits long, and the nonce argument must be 96 bits long. The nonce must be
// randomly generated or used only once. This Stream instance must not be used
// to encrypt more than 2^38 bytes (256 gigabytes).
func NewRFC(key, nonce []byte) (cipher.Stream, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKey
	}

	if len(nonce) != RFCNonceSize {
		return nil, ErrInvalidNonce
	}

	return ref.NewRFC(key, nonce)
}

// NewDraft creates and returns a new cipher.Stream. The key argument must be
// 256 bits long, and the nonce argument must be 64 bits long. The nonce must
// be randomly generated or used only once. This Stream instance must not be
// used to encrypt more than 2^70 bytes (~1 zettabyte).
func NewDraft(key, nonce []byte) (cipher.Stream, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKey
	}

	if len(nonce) != DraftNonceSize {
		return nil, ErrInvalidNonce
	}

	return ref.NewDraft(key, nonce)
}

// NewXChaCha creates and returns a new cipher.Stream. The key argument must be
// 256 bits long, and the nonce argument must be 192 bits long. The nonce must
// be randomly generated or only used once. This Stream instance must not be
// used to encrypt more than 2^70 bytes (~1 zetta byte).
func NewXChaCha(key, nonce []byte) (cipher.Stream, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKey
	}

	if len(nonce) != XNonceSize {
		return nil, ErrInvalidNonce
	}

	return ref.NewXChaCha(key, nonce)
}
