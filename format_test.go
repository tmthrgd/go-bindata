// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/tools/imports"
)

func TestFormatting(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c := &Config{Package: "main"}
			test.config(c)

			var buf bytes.Buffer
			if err := testGenerate(&buf, c); err != nil {
				t.Fatal(err)
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
