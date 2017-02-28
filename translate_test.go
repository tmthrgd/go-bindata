// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"
)

var testCases = [...]struct {
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
	{"prefix", func(c *Config) { c.Prefix = "testdata" }},
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
}

var testPaths = [...]struct {
	path      string
	recursive bool
}{
	{"testdata", true},
	{"LICENSE", false},
	{"README.md", false},
}

var gencode = flag.String("gencode", "", "write generated code to specified directory")

func testGenerate(w io.Writer, c *Config) error {
	g, err := New(c)
	if err != nil {
		return err
	}

	for _, path := range testPaths {
		if err = g.FindFiles(path.path, path.recursive); err != nil {
			return err
		}
	}

	_, err = g.WriteTo(w)
	return err
}

func TestGenerate(t *testing.T) {
	if *gencode == "" {
		t.Skip("skipping test as -gencode flag not provided")
	}

	if err := os.Mkdir(*gencode, 0777); err != nil && !os.IsExist(err) {
		t.Fatal(err)
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c := &Config{Package: "main"}
			test.config(c)

			f, err := os.Create(filepath.Join(*gencode, test.name+".go"))
			if err != nil {
				t.Fatal(err)
			}

			err = testGenerate(f, c)
			f.Close()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
