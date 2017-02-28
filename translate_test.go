// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"flag"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

var testCases = [...]struct {
	name string
	opts func(*GenerateOptions)
}{
	{"default", func(*GenerateOptions) {}},
	{"old-default", func(o *GenerateOptions) {
		o.Package = "main"
		o.MemCopy = true
		o.Compress = true
		o.Metadata = true
		// The AssetDir API currently produces
		// wrongly formatted code. We're going
		// to skip it for now.
		/*o.AssetDir = true
		o.Restore = true*/
		o.DecompressOnce = true
	}},
	{"debug", func(o *GenerateOptions) { o.Debug = true }},
	{"dev", func(o *GenerateOptions) { o.Dev = true }},
	{"tags", func(o *GenerateOptions) { o.Tags = "!x" }},
	{"package", func(o *GenerateOptions) { o.Package = "test" }},
	{"compress", func(o *GenerateOptions) { o.Compress = true }},
	{"copy", func(o *GenerateOptions) { o.MemCopy = true }},
	{"metadata", func(o *GenerateOptions) { o.Metadata = true }},
	{"decompress-once", func(o *GenerateOptions) {
		o.Compress = true
		o.DecompressOnce = true
	}},
	{"hash-dir", func(o *GenerateOptions) { o.HashFormat = DirHash }},
	{"hash-suffix", func(o *GenerateOptions) { o.HashFormat = NameHashSuffix }},
	{"hash-hashext", func(o *GenerateOptions) { o.HashFormat = HashWithExt }},
	{"hash-unchanged", func(o *GenerateOptions) { o.HashFormat = NameUnchanged }},
	{"hash-enc-b32", func(o *GenerateOptions) {
		o.HashEncoding = Base32Hash
		o.HashFormat = DirHash
	}},
	{"hash-enc-b64", func(o *GenerateOptions) {
		o.HashEncoding = Base64Hash
		o.HashFormat = DirHash
	}},
	{"hash-key", func(o *GenerateOptions) {
		o.HashKey = []byte{0x00, 0x11, 0x22, 0x33}
		o.HashFormat = DirHash
	}},
}

var testPaths = map[string]*FindFilesOptions{
	"testdata":        {Recursive: true},
	"testdata/ogqS":   {Prefix: "testdata"},
	"CONTRIBUTING.md": nil,
	"LICENSE":         nil,
	"README.md":       nil,
}

var gencode = flag.String("gencode", "", "write generated code to specified directory")

func testFiles() (Files, error) {
	var all Files

	for path, opts := range testPaths {
		files, err := FindFiles(path, opts)
		if err != nil {
			return nil, err
		}

		all = append(all, files...)
	}

	sort.Sort(all)
	return all, nil
}

func TestGenerate(t *testing.T) {
	if *gencode == "" {
		t.Skip("skipping test as -gencode flag not provided")
	}

	files, err := testFiles()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir(*gencode, 0777); err != nil && !os.IsExist(err) {
		t.Fatal(err)
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			o := &GenerateOptions{Package: "main"}
			test.opts(o)

			f, err := os.Create(filepath.Join(*gencode, test.name+".go"))
			if err != nil {
				t.Fatal(err)
			}

			err = files.Generate(f, o)
			f.Close()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
