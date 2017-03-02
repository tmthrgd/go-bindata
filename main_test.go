// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func TestMain(m *testing.M) {
	flag.Parse()

	setupTestCases()
	setupTestFiles()

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
