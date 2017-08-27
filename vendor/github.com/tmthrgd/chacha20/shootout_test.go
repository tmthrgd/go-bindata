// Copyright 2014 Coda Hale. All rights reserved.
// Use of this source code is governed by an MIT
// License that can be found in the LICENSE file.

package chacha20

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rc4"
	"testing"

	codahale "github.com/codahale/chacha20"
	ref "github.com/tmthrgd/chacha20/internal/ref"
)

type size struct {
	name string
	l    int
}

var sizes = []size{
	{"32", 32},
	{"128", 128},
	{"1K", 1 * 1024},
	{"16K", 16 * 1024},
	{"128K", 128 * 1024},
	{"1M", 1024 * 1024},
}

func benchmarkStream(b *testing.B, c cipher.Stream, l int) {
	input := make([]byte, l)
	output := make([]byte, l)

	b.SetBytes(int64(l))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.XORKeyStream(output, input)
	}
}

func BenchmarkChaCha20Codahale(b *testing.B) {
	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			key := make([]byte, codahale.KeySize)
			nonce := make([]byte, codahale.NonceSize)
			c, _ := codahale.New(key, nonce)

			benchmarkStream(b, c, size.l)
		})
	}
}

func BenchmarkChaCha20Go(b *testing.B) {
	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			key := make([]byte, KeySize)
			nonce := make([]byte, RFCNonceSize)
			c, _ := ref.NewRFC(key, nonce)

			benchmarkStream(b, c, size.l)
		})
	}
}

func BenchmarkChaCha20x64(b *testing.B) {
	if useRef {
		b.Skip("skipping: do not have x64 implementation")
	}

	oldAVX, oldAVX2 := useAVX, useAVX2
	useAVX, useAVX2 = false, false
	defer func() {
		useAVX, useAVX2 = oldAVX, oldAVX2
	}()

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			key := make([]byte, KeySize)
			nonce := make([]byte, RFCNonceSize)
			c, _ := NewRFC(key, nonce)

			benchmarkStream(b, c, size.l)
		})
	}
}

func BenchmarkChaCha20AVX(b *testing.B) {
	if !useAVX {
		b.Skip("skipping: do not have AVX implementation")
	}

	oldAVX, oldAVX2 := useAVX, useAVX2
	useAVX, useAVX2 = true, false
	defer func() {
		useAVX, useAVX2 = oldAVX, oldAVX2
	}()

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			key := make([]byte, KeySize)
			nonce := make([]byte, RFCNonceSize)
			c, _ := NewRFC(key, nonce)

			benchmarkStream(b, c, size.l)
		})
	}
}

func BenchmarkChaCha20AVX2(b *testing.B) {
	if !useAVX2 {
		b.Skip("skipping: do not have AVX2 implementation")
	}

	oldAVX, oldAVX2 := useAVX, useAVX2
	useAVX, useAVX2 = false, true
	defer func() {
		useAVX, useAVX2 = oldAVX, oldAVX2
	}()

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			key := make([]byte, KeySize)
			nonce := make([]byte, RFCNonceSize)
			c, _ := NewRFC(key, nonce)

			benchmarkStream(b, c, size.l)
		})
	}
}

func BenchmarkAESCTR(b *testing.B) {
	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			key := make([]byte, 32)
			a, _ := aes.NewCipher(key)

			iv := make([]byte, aes.BlockSize)
			c := cipher.NewCTR(a, iv)

			benchmarkStream(b, c, size.l)
		})
	}
}

func BenchmarkAESGCM(b *testing.B) {
	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			key := make([]byte, 32)
			a, _ := aes.NewCipher(key)
			c, _ := cipher.NewGCM(a)

			nonce := make([]byte, c.NonceSize())

			input := make([]byte, size.l)
			output := make([]byte, 0, size.l+c.Overhead())

			b.SetBytes(int64(size.l))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				c.Seal(output, nonce, input, nil)
			}
		})
	}
}

func BenchmarkRC4(b *testing.B) {
	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			key := make([]byte, 32)
			c, _ := rc4.NewCipher(key)

			benchmarkStream(b, c, size.l)
		})
	}
}
