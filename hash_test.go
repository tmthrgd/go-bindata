// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"os"
	"testing"
)

func BenchmarkHashFile(b *testing.B) {
	files, err := testFiles()
	if err != nil {
		b.Fatal(err)
	}

	file, size := files[0], int64(0)
	for _, f := range files {
		info, err := os.Stat(file.Path)
		if err != nil {
			b.Error(err)
			continue
		}

		if l := info.Size(); l > size {
			file = f
			size = l
		}
	}

	first := true
	for _, enc := range [...]HashEncoding{HexHash, Base32Hash, Base64Hash} {
		b.Run(enc.String(), func(b *testing.B) {
			for _, fmt := range [...]HashFormat{NameUnchanged, DirHash, NameHashSuffix, HashWithExt} {
				b.Run(fmt.String(), func(b *testing.B) {
					opts := &GenerateOptions{
						HashFormat:   fmt,
						HashEncoding: enc,
						HashLength:   32,
					}

					if first {
						b.Logf("hashing file %s of size %dB", file.Path, size)
						first = false
					}

					for n := 0; n < b.N; n++ {
						asset := binAsset{
							Path: file.Path,
							Name: file.Name,
						}

						if err := hashFile(opts, &asset); err != nil {
							b.Fatal(err)
						}

						if asset.Hash == nil || (fmt != NameUnchanged && asset.Name == file.Name) {
							b.Fatal("hashFile failed")
						}
					}
				})
			}
		})
	}
}
