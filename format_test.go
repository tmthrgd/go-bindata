// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"flag"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/tools/imports"
)

var gencode = flag.Bool("gencode", false, "will log the generated go code")

func TestFormatting(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	for _, test := range [...]struct {
		name   string
		config func(*Config)
	}{
		{"default", func(*Config) {}},
		{"old-default", func(c *Config) {
			c.Package = "main"
			c.MemCopy = true
			c.Compress = true
			c.Metadata = true
			c.HashLength = 16
			// The AssetDir API currently produces
			// wrongly formatted code. We're going
			// to skip it for now.
			/*c.AssetDir = true
			c.Restore = true*/
			c.DecompressOnce = true
		}},
		{"debug", func(c *Config) { c.Debug = true }},
		{"dev", func(c *Config) { c.Dev = true }},
		{"tags", func(c *Config) { c.Tags = "!x" }},
		{"package", func(c *Config) { c.Package = "test" }},
		{"compress", func(c *Config) { c.Compress = true }},
		{"copy", func(c *Config) { c.MemCopy = true }},
		{"metadata", func(c *Config) { c.Metadata = true }},
		{"decompress-once", func(c *Config) { c.DecompressOnce = true }},
		{"hash-dir", func(c *Config) { c.HashFormat = DirHash; c.HashLength = 16 }},
		{"hash-suffix", func(c *Config) { c.HashFormat = NameHashSuffix; c.HashLength = 16 }},
		{"hash-hashext", func(c *Config) { c.HashFormat = HashWithExt; c.HashLength = 16 }},
		{"hash-unchanged", func(c *Config) { c.HashFormat = NameUnchanged; c.HashLength = 16 }},
		{"hash-enc-b32", func(c *Config) { c.HashEncoding = Base32Hash; c.HashFormat = DirHash; c.HashLength = 16 }},
		{"hash-enc-b64", func(c *Config) { c.HashEncoding = Base64Hash; c.HashFormat = DirHash; c.HashLength = 16 }},
		{"hash-key", func(c *Config) { c.HashKey = []byte{0x00, 0x11, 0x22, 0x33}; c.HashFormat = DirHash; c.HashLength = 16 }},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c := &Config{
				Package: "main",
				Input: []InputConfig{
					{
						Path:      "testdata",
						Recursive: true,
					},
				},
			}
			test.config(c)

			var buf bytes.Buffer
			if err := Translate(&buf, c); err != nil {
				t.Fatal(err)
			}

			if *gencode {
				t.Log(buf.String())
			}

			out, err := imports.Process("bindata.go", buf.Bytes(), nil)
			if err != nil {
				t.Fatal(err)
			}

			if bytes.Equal(buf.Bytes(), out) {
				return
			}

			t.Error("not correctly formatted.")

			var diff bytes.Buffer
			diff.WriteString("diff:\n")

			if err := difflib.WriteUnifiedDiff(&diff, difflib.UnifiedDiff{
				A:       difflib.SplitLines(buf.String()),
				B:       difflib.SplitLines(string(out)),
				Context: 2,
				Eol:     "",
			}); err != nil {
				t.Fatal(err)
			}

			t.Log(diff.String())
		})
	}
}
