// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"testing/quick"
	"unicode"
	"unicode/utf8"
)

type testCase struct {
	name string
	opts func(*GenerateOptions)
}

var testCases = []testCase{
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

var (
	gencode = flag.String("gencode", "", "write generated code to specified directory")

	randTestCases = flag.Uint("randtests", 25, "")
)

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
		vo.Package = identifier(vo.Package)
		vo.Mode &= os.ModePerm
		vo.HashFormat = HashFormat(int(uint(vo.HashFormat) % uint(HashWithExt+1)))
		vo.HashEncoding = HashEncoding(int(uint(vo.HashEncoding) % uint(Base64Hash+1)))
		// The AssetDir API currently produces
		// wrongly formatted code. We're going
		// to skip it for now.
		vo.AssetDir = false
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

		testCases = append(testCases, testCase{
			fmt.Sprintf("random-#%d", i+1),
			func(o *GenerateOptions) { *o = *vo },
		})
	}

	os.Exit(m.Run())
}

// identifier removes all characters from a string that are not valid in
// an identifier according to the Go Programming Language Specification.
//
// The logic in the switch statement was taken from go/source package:
// https://github.com/golang/go/blob/a1a688fa0012f7ce3a37e9ac0070461fe8e3f28e/src/go/scanner/scanner.go#L257-#L271
func identifier(val string) string {
	return strings.TrimLeftFunc(strings.Map(func(ch rune) rune {
		switch {
		case 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' ||
			ch >= utf8.RuneSelf && unicode.IsLetter(ch):
			return ch
		case '0' <= ch && ch <= '9' ||
			ch >= utf8.RuneSelf && unicode.IsDigit(ch):
			return ch
		default:
			return -1
		}
	}, val), unicode.IsDigit)
}

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
