// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"flag"
	"os"
	"sort"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

var testPaths = map[string]*FindFilesOptions{
	"testdata":               {Recursive: true},
	"testdata/ab6.bin":       {Prefix: "testdata"},
	"testdata/ogqS":          {Prefix: "testdata"},
	"testdata/ogqS/qsDM.bin": {Prefix: "testdata/ogqS"},
}

var (
	testFiles    Files
	testFilesErr error
)

func TestMain(m *testing.M) {
	flag.Parse()

	setupTestCases()

	paths := make([]string, 0, len(testPaths))
	for path := range testPaths {
		paths = append(paths, path)
	}

	sort.Strings(paths)

	for _, path := range paths {
		files, err := FindFiles(path, testPaths[path])
		if err != nil {
			testFilesErr = err
			break
		}

		testFiles = append(testFiles, files...)
	}

	os.Exit(m.Run())
}

func testDiff(a, b string) (string, error) {
	var diff bytes.Buffer
	diff.WriteString("diff:\n")

	if err := difflib.WriteUnifiedDiff(&diff, difflib.UnifiedDiff{
		A:       difflib.SplitLines(a),
		B:       difflib.SplitLines(b),
		Context: 3,
		Eol:     "",
	}); err != nil {
		return "", nil
	}

	return diff.String(), nil
}
