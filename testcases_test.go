// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"reflect"

	"github.com/tmthrgd/go-bindata/internal/identifier"
	"golang.org/x/crypto/blake2b"
)

var testCases = map[string]func(*GenerateOptions){
	"default": func(*GenerateOptions) {},
	"old-default": func(o *GenerateOptions) {
		*o = GenerateOptions{
			Package:        "main",
			MemCopy:        true,
			Compress:       true,
			Metadata:       true,
			AssetDir:       true,
			Restore:        true,
			DecompressOnce: true,
		}
	},
	"debug":    func(o *GenerateOptions) { o.Debug = true },
	"dev":      func(o *GenerateOptions) { o.Dev = true },
	"tags":     func(o *GenerateOptions) { o.Tags = "!x" },
	"package":  func(o *GenerateOptions) { o.Package = "test" },
	"compress": func(o *GenerateOptions) { o.Compress = true },
	"copy":     func(o *GenerateOptions) { o.MemCopy = true },
	"metadata": func(o *GenerateOptions) { o.Metadata = true },
	"decompress-once": func(o *GenerateOptions) {
		o.Compress = true
		o.DecompressOnce = true
	},
	"hash-unchanged": func(o *GenerateOptions) {
		o.Hash, _ = blake2b.New512(nil)
	},
	"hash-dir": func(o *GenerateOptions) {
		o.Hash, _ = blake2b.New512(nil)
		o.HashFormat = DirHash
	},
	"hash-suffix": func(o *GenerateOptions) {
		o.Hash, _ = blake2b.New512(nil)
		o.HashFormat = NameHashSuffix
	},
	"hash-hashext": func(o *GenerateOptions) {
		o.Hash, _ = blake2b.New512(nil)
		o.HashFormat = HashWithExt
	},
	"hash-enc-b32": func(o *GenerateOptions) {
		o.Hash, _ = blake2b.New512(nil)
		o.HashEncoding = Base32Hash
		o.HashFormat = DirHash
	},
	"hash-enc-b64": func(o *GenerateOptions) {
		o.Hash, _ = blake2b.New512(nil)
		o.HashEncoding = Base64Hash
		o.HashFormat = DirHash
	},
	"hash-copy": func(o *GenerateOptions) {
		o.MemCopy = true
		o.Hash, _ = blake2b.New512(nil)
	},
	"asset-dir": func(o *GenerateOptions) { o.AssetDir = true },
}

var randTestCases = flag.Uint("randtests", 25, "the number of random test cases to add")

func setupTestCases() {
	t := reflect.TypeOf(GenerateOptions{})

	for i := uint(0); i < *randTestCases; i++ {
		rand := rand.New(rand.NewSource(int64(i)))

		v, ok := sizedValue(t, rand, complexSize)
		if !ok {
			panic("sizedValue failed")
		}

		vo := v.Addr().Interface().(*GenerateOptions)
		vo.Package = identifier.Identifier(vo.Package)
		vo.Mode &= os.ModePerm
		vo.Metadata = vo.Metadata && (vo.Mode == 0 || vo.ModTime == 0)
		vo.HashFormat = HashFormat(int(uint(vo.HashFormat) % uint(HashWithExt+1)))
		vo.HashEncoding = HashEncoding(int(uint(vo.HashEncoding) % uint(Base64Hash+1)))
		vo.Restore = vo.Restore && vo.AssetDir

		if vo.Package == "" {
			vo.Package = "main"
		}

		if vo.Debug || vo.Dev {
			vo.Hash = nil
		}

		testCases[fmt.Sprintf("random-#%d", i+1)] = func(o *GenerateOptions) { *o = *vo }
	}
}
