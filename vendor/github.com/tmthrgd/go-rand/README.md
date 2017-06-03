# go-rand

[![GoDoc](https://godoc.org/github.com/tmthrgd/go-rand?status.svg)](https://godoc.org/github.com/tmthrgd/go-rand)
[![Build Status](https://travis-ci.org/tmthrgd/go-rand.svg?branch=master)](https://travis-ci.org/tmthrgd/go-rand)

A cryptographically secure pseudo-random number generator built from ChaCha20.

**There are no security guarantees offered.**

## Benchmark

```
BenchmarkCryptoRand-8   	      20	  63222315 ns/op	  16.59 MB/s	[crypto/rand]
BenchmarkMathRand-8     	     500	   2572662 ns/op	 407.58 MB/s	[math/rand]
BenchmarkReader-8       	    2000	    732662 ns/op	1431.19 MB/s	[tmthrgd/go-rand - AVX only]
```

## License

Unless otherwise noted, the go-rand source files are distributed under the Modified BSD License found in the LICENSE file.
