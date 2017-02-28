// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import "testing"

func BenchmarkFindFiles(b *testing.B) {
	for _, test := range testCases {
		test := test
		b.Run(test.name, func(b *testing.B) {
			c := &Config{Package: "main"}
			test.config(c)

			g, err := New(c)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				if err = g.FindFiles("testdata", &FindFilesOptions{
					Prefix:    "testdata",
					Recursive: true,
				}); err != nil {
					b.Fatal(err)
				}

				g.toc = nil
			}
		})
	}
}

func BenchmarkWriteTo(b *testing.B) {
	for _, test := range testCases {
		test := test
		b.Run(test.name, func(b *testing.B) {
			c := &Config{Package: "main"}
			test.config(c)

			g, err := New(c)
			if err != nil {
				b.Fatal(err)
			}

			for path, opts := range testPaths {
				if err = g.FindFiles(path, opts); err != nil {
					b.Fatal(err)
				}
			}

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				if _, err = g.WriteTo(nopWriter{}); err != nil {
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
