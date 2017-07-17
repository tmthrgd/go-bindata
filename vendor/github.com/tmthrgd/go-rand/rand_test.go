// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.
//
// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rand

import (
	"bytes"
	"compress/flate"
	"errors"
	"io"
	"testing"
	"testing/quick"
)

var testVectors = []struct {
	seed   []byte
	expect []byte
}{
	{
		[]byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
		[]byte{
			0x76, 0xb8, 0xe0, 0xad, 0xa0, 0xf1, 0x3d, 0x90, 0x40, 0x5d, 0x6a, 0xe5, 0x53, 0x86, 0xbd, 0x28,
			0xbd, 0xd2, 0x19, 0xb8, 0xa0, 0x8d, 0xed, 0x1a, 0xa8, 0x36, 0xef, 0xcc, 0x8b, 0x77, 0x0d, 0xc7,
			0xda, 0x41, 0x59, 0x7c, 0x51, 0x57, 0x48, 0x8d, 0x77, 0x24, 0xe0, 0x3f, 0xb8, 0xd8, 0x4a, 0x37,
			0x6a, 0x43, 0xb8, 0xf4, 0x15, 0x18, 0xa1, 0x1c, 0xc3, 0x87, 0xb6, 0x69, 0xb2, 0xee, 0x65, 0x86,
			0x9f, 0x07, 0xe7, 0xbe, 0x55, 0x51, 0x38, 0x7a, 0x98, 0xba, 0x97, 0x7c, 0x73, 0x2d, 0x08, 0x0d,
			0xcb, 0x0f, 0x29, 0xa0, 0x48, 0xe3, 0x65, 0x69, 0x12, 0xc6, 0x53, 0x3e, 0x32, 0xee, 0x7a, 0xed,
			0x29, 0xb7, 0x21, 0x76, 0x9c, 0xe6, 0x4e, 0x43, 0xd5, 0x71, 0x33, 0xb0, 0x74, 0xd8, 0x39, 0xd5,
			0x31, 0xed, 0x1f, 0x28, 0x51, 0x0a, 0xfb, 0x45, 0xac, 0xe1, 0x0a, 0x1f, 0x4b, 0x79, 0x4d, 0x6f,
		},
	},
}

func TestReadFull(t *testing.T) {
	t.Parallel()

	var seed [SeedSize]byte

	r, err := New(seed[:])
	if err != nil {
		t.Fatal(err)
	}

	var scratch, zero [128 * 1024 * 1024]byte

	n, err := r.Read(scratch[:])
	if err != nil {
		t.Error(err)
	}

	if n != len(scratch) {
		t.Errorf("expected read to return %d bytes, got %d", len(scratch), n)
	}

	if bytes.Equal(scratch[:], zero[:]) {
		t.Error("read failed")
	}
}

func TestVectors(t *testing.T) {
	t.Parallel()

	for i, vector := range testVectors {
		t.Logf("running test vector %d\n", i)

		r, err := New(vector.seed)
		if err != nil {
			t.Error(err)
			continue
		}

		data := make([]byte, len(vector.expect))

		if _, err := r.Read(data); err != nil {
			t.Error(err)
			continue
		}

		if !bytes.Equal(data, vector.expect) {
			t.Error("invalid output")
			t.Logf("\texpected %x\n", vector.expect)
			t.Logf("\tgot      %x\n", data)
		}
	}
}

func TestNewWithSeed(t *testing.T) {
	t.Parallel()

	var seed [SeedSize]byte

	if _, err := New(seed[:]); err != nil {
		t.Error(err)
	}
}

func TestNewNoSeed(t *testing.T) {
	t.Parallel()

	if _, err := New(nil); err != nil {
		t.Error(err)
	}
}

func TestInvalidSeed(t *testing.T) {
	t.Parallel()

	var seed [SeedSize - 1]byte

	if _, err := New(seed[:]); err == nil {
		t.Error("expected error, got <nil>")
	}
}

func TestGlobalRead(t *testing.T) {
	t.Parallel()

	var scratch, zero [157]byte

	n, err := Read(scratch[:])
	if err != nil {
		t.Error(err)
	}

	if n != len(scratch) {
		t.Errorf("expected read to return %d bytes, got %d", len(scratch), n)
	}

	if bytes.Equal(scratch[:], zero[:]) {
		t.Error("read failed")
	}
}

func TestRead(t *testing.T) {
	t.Parallel()

	var seed [SeedSize]byte

	r, err := New(seed[:])
	if err != nil {
		t.Fatal(err)
	}

	var scratch, zero [157]byte

	n, err := r.Read(scratch[:])
	if err != nil {
		t.Error(err)
	}

	if n != len(scratch) {
		t.Errorf("expected read to return %d bytes, got %d", len(scratch), n)
	}

	if bytes.Equal(scratch[:], zero[:]) {
		t.Error("read failed")
	}
}

func TestLongRead(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping: short test flag")
	}

	var seed [SeedSize]byte

	r, err := New(seed[:])
	if err != nil {
		t.Fatal(err)
	}

	var scratch [1024 * 1024 * 1024]byte

	n, err := r.Read(scratch[:])
	if err != nil {
		t.Error(err)
	}

	if n != len(scratch) {
		t.Errorf("expected read to return %d bytes, got %d", len(scratch), n)
	}

	for _, v := range scratch[:] {
		if v != 0 {
			return
		}
	}

	t.Error("read failed")
}

func testReseed(t *testing.T, partway bool) {
	var seed [SeedSize]byte

	r, err := New(seed[:])
	if err != nil {
		t.Fatal(err)
	}

	var scratch, zero [194]byte

	if partway {
		r.(*reader).budget = 97
	} else {
		r.(*reader).budget = 0
	}

	n, err := r.Read(scratch[:])
	if err != nil {
		t.Error(err)
	}

	if n != len(scratch) {
		t.Errorf("expected read to return %d bytes, got %d", len(scratch), n)
	}

	if bytes.Equal(scratch[:], zero[:]) {
		t.Error("read failed")
	}
}

func TestReseed(t *testing.T) {
	t.Parallel()

	testReseed(t, false)
}

func TestReseedPartway(t *testing.T) {
	t.Parallel()

	testReseed(t, true)
}

func TestReadEmpty(t *testing.T) {
	t.Parallel()

	var seed [SeedSize]byte

	r, err := New(seed[:])
	if err != nil {
		t.Fatal(err)
	}

	if n, err := r.Read(make([]byte, 0)); n != 0 || err != nil {
		t.Fatalf("Read(make([]byte, 0)) = %d, %v", n, err)
	}

	if n, err := r.Read(nil); n != 0 || err != nil {
		t.Fatalf("Read(nil) = %d, %v", n, err)
	}
}

func TestCompressability(t *testing.T) {
	t.Parallel()

	var seed [SeedSize]byte

	r, err := New(seed[:])
	if err != nil {
		t.Fatal(err)
	}

	var n int = 3e7
	if testing.Short() {
		n = 1e5
	}

	var z bytes.Buffer

	f, err := flate.NewWriter(&z, flate.DefaultCompression)
	if err != nil {
		t.Fatal(err)
	}

	if nn, err := io.CopyN(f, r, int64(n)); nn != int64(n) || err != nil {
		t.Fatalf("io.CopyN(f, r, n) = %d, %s", nn, err)
	}

	f.Close()

	if z.Len() < n*999/1000 {
		t.Errorf("compressed %d -> %d: %f%%", n, z.Len(), float64(z.Len())/float64(n)*100)
	} else {
		t.Logf("compressed %d -> %d: %f%%", n, z.Len(), float64(z.Len())/float64(n)*100)
	}
}

func TestSeededCompressability(t *testing.T) {
	t.Parallel()

	r, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}

	var n int = 3e7
	if testing.Short() {
		n = 1e5
	}

	var z bytes.Buffer

	f, err := flate.NewWriter(&z, flate.DefaultCompression)
	if err != nil {
		t.Fatal(err)
	}

	if nn, err := io.CopyN(f, r, int64(n)); nn != int64(n) || err != nil {
		t.Fatalf("io.CopyN(f, r, n) = %d, %s", nn, err)
	}

	f.Close()

	if z.Len() < n*999/1000 {
		t.Errorf("compressed %d -> %d: %f%%", n, z.Len(), float64(z.Len())/float64(n)*100)
	} else {
		t.Logf("compressed %d -> %d: %f%%", n, z.Len(), float64(z.Len())/float64(n)*100)
	}
}

func TestIndependence(t *testing.T) {
	t.Parallel()

	var seed [SeedSize]byte

	a, err := New(seed[:])
	if err != nil {
		t.Fatal(err)
	}

	seed[0] = 1

	b, err := New(seed[:])
	if err != nil {
		t.Fatal(err)
	}

	var n int = 3e7
	if testing.Short() {
		n = 1e5
	}

	n /= 2

	var z bytes.Buffer

	f, err := flate.NewWriter(&z, flate.DefaultCompression)
	if err != nil {
		t.Fatal(err)
	}

	if nn, err := io.CopyN(f, a, int64(n)); nn != int64(n) || err != nil {
		t.Fatalf("io.CopyN(f, a, n) = %d, %s", nn, err)
	}

	if nn, err := io.CopyN(f, b, int64(n)); nn != int64(n) || err != nil {
		t.Fatalf("io.CopyN(f, b, n) = %d, %s", nn, err)
	}

	f.Close()

	if z.Len() < 2*n*999/1000 {
		t.Errorf("compressed %d -> %d: %f%%", 2*n, z.Len(), float64(z.Len())/float64(2*n)*100)
	} else {
		t.Logf("compressed %d -> %d: %f%%", 2*n, z.Len(), float64(z.Len())/float64(2*n)*100)
	}
}

func TestReseedIndependence(t *testing.T) {
	t.Parallel()

	var seed [SeedSize]byte

	r, err := New(seed[:])
	if err != nil {
		t.Fatal(err)
	}

	var n, m int = 3e7, 5
	if testing.Short() {
		n, m = 1e5, 4
	}

	n /= m

	var z bytes.Buffer

	f, err := flate.NewWriter(&z, flate.DefaultCompression)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < m; i++ {
		if nn, err := io.CopyN(f, r, int64(n)); nn != int64(n) || err != nil {
			t.Fatalf("io.CopyN(f, r, n) = %d, %s", nn, err)
		}

		r.(*reader).budget = 0
	}

	f.Close()

	if z.Len() < m*n*999/1000 {
		t.Errorf("compressed %d -> %d: %f%%", m*n, z.Len(), float64(z.Len())/float64(m*n)*100)
	} else {
		t.Logf("compressed %d -> %d: %f%%", m*n, z.Len(), float64(z.Len())/float64(m*n)*100)
	}
}

func TestDeterministic(t *testing.T) {
	t.Parallel()

	var seed [SeedSize]byte

	a, err := New(seed[:])
	if err != nil {
		t.Fatal(err)
	}

	b, err := New(seed[:])
	if err != nil {
		t.Fatal(err)
	}

	if err = quick.CheckEqual(func() ([]byte, error) {
		dst := make([]byte, 1024)

		n, err := a.Read(dst)
		if err != nil {
			return nil, err
		}

		return dst[:n], nil
	}, func() ([]byte, error) {
		dst := make([]byte, 1024)

		n, err := b.Read(dst)
		if err != nil {
			return nil, err
		}

		return dst[:n], nil
	}, &quick.Config{
		MaxCountScale: 100,
	}); err != nil {
		t.Error(err)
	}
}

func TestSeedOnFirstRead(t *testing.T) {
	r := &seedOnFirstRead{r: nil}

	r.seed()
	if r.r == nil {
		t.Errorf("(*seedOnFirstRead).seed did not create reader")
	}
	if r.err != nil {
		t.Error(r.err)
	}

	var scratch [1]byte

	rr := r.r
	if _, err := r.Read(scratch[:]); err != nil {
		t.Error(err)
	}

	if r.r == rr {
		t.Error("(*seedOnFirstRead).Read did not create reader")
	}

	rr = r.r
	if _, err := r.Read(scratch[:]); err != nil {
		t.Error(err)
	}

	if r.r != rr {
		t.Error("(*seedOnFirstRead).Read seeded twice")
	}

	terr := errors.New("test error")
	r.err = terr

	if _, err := r.Read(scratch[:]); err != terr {
		t.Errorf("(*seedOnFirstRead).Read did not return error, expected %v, got %v", terr, err)
	}
}
