// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package rand

import (
	crand "crypto/rand"
	"io"
	mrand "math/rand"
	"testing"
)

const benchSize = 1024 * 1024

func benchmarkReader(b *testing.B, r io.Reader) {
	b.SetBytes(benchSize)

	output := make([]byte, benchSize)

	for i := 0; i < b.N; i++ {
		io.ReadFull(r, output)
	}
}

func BenchmarkCryptoRand(b *testing.B) {
	benchmarkReader(b, crand.Reader)
}

func BenchmarkMathRand(b *testing.B) {
	r := mrand.New(mrand.NewSource(1))

	benchmarkReader(b, r)
}

func BenchmarkReader(b *testing.B) {
	benchmarkReader(b, Reader)
}
