// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"testing"
	"testing/quick"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/tmthrgd/go-bindata/internal/identifier"
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
	"hash-dir":       func(o *GenerateOptions) { o.HashFormat = DirHash },
	"hash-suffix":    func(o *GenerateOptions) { o.HashFormat = NameHashSuffix },
	"hash-hashext":   func(o *GenerateOptions) { o.HashFormat = HashWithExt },
	"hash-unchanged": func(o *GenerateOptions) { o.HashFormat = NameUnchanged },
	"hash-enc-b32": func(o *GenerateOptions) {
		o.HashEncoding = Base32Hash
		o.HashFormat = DirHash
	},
	"hash-enc-b64": func(o *GenerateOptions) {
		o.HashEncoding = Base64Hash
		o.HashFormat = DirHash
	},
	"hash-key": func(o *GenerateOptions) {
		o.HashKey = []byte{0x00, 0x11, 0x22, 0x33}
		o.HashFormat = DirHash
	},
}

var testPaths = map[string]*FindFilesOptions{
	"testdata":        {Recursive: true},
	"testdata/ogqS":   {Prefix: "testdata"},
	"CONTRIBUTING.md": nil,
	"LICENSE":         nil,
	"README.md":       nil,
}

var (
	testFiles    Files
	testFilesErr error
)

var randTestCases = flag.Uint("randtests", 25, "")

func TestMain(m *testing.M) {
	flag.Parse()

	t := reflect.TypeOf(GenerateOptions{})

	for i := uint(0); i < *randTestCases; i++ {
		rand := rand.New(rand.NewSource(int64(i)))

		v, ok := quick.Value(t, rand)
		if !ok {
			panic("quick.Value failed")
		}

		vo := v.Addr().Interface().(*GenerateOptions)
		vo.Package = identifier.Identifier(vo.Package)
		vo.Mode &= os.ModePerm
		vo.HashFormat = HashFormat(int(uint(vo.HashFormat) % uint(HashWithExt+1)))
		vo.HashEncoding = HashEncoding(int(uint(vo.HashEncoding) % uint(Base64Hash+1)))
		vo.Restore = vo.Restore && vo.AssetDir

		if vo.Package == "" {
			vo.Package = "main"
		}

		switch vo.HashEncoding {
		case HexHash:
			vo.HashLength %= maxHexLength
		case Base32Hash:
			vo.HashLength %= maxB32Length
		case Base64Hash:
			vo.HashLength %= maxB64Length
		}

		if vo.Debug || vo.Dev {
			vo.HashFormat = NoHash
		}

		testCases[fmt.Sprintf("random-#%d", i+1)] = func(o *GenerateOptions) { *o = *vo }
	}

	for path, opts := range testPaths {
		files, err := FindFiles(path, opts)
		if err != nil {
			testFilesErr = err
			break
		}

		testFiles = append(testFiles, files...)
	}

	sort.Sort(testFiles)

	os.Exit(m.Run())
}

func testDiff(a, b string) (string, error) {
	var diff bytes.Buffer
	diff.WriteString("diff:\n")

	if err := difflib.WriteUnifiedDiff(&diff, difflib.UnifiedDiff{
		A:       difflib.SplitLines(a),
		B:       difflib.SplitLines(b),
		Context: 2,
		Eol:     "",
	}); err != nil {
		return "", nil
	}

	return diff.String(), nil
}
