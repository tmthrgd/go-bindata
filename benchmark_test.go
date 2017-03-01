// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import "testing"

func BenchmarkFindFiles(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if _, err := FindFiles("testdata", &FindFilesOptions{
			Prefix:    "testdata",
			Recursive: true,
		}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerate(b *testing.B) {
	if testFilesErr != nil {
		b.Fatal(testFilesErr)
	}

	for name, opts := range testCases {
		name, opts := name, opts
		b.Run(name, func(b *testing.B) {
			o := &GenerateOptions{Package: "main"}
			opts(o)

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				if err := testFiles.Generate(nopWriter{}, o); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

type nopWriter struct{}

func (nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
