// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"testing"

	"golang.org/x/tools/imports"
)

func TestFormatting(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	files, err := testFiles()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			o := &GenerateOptions{Package: "main"}
			test.opts(o)

			var buf bytes.Buffer
			if err := files.Generate(&buf, o); err != nil {
				t.Fatal(err)
			}

			out, err := imports.Process("bindata.go", buf.Bytes(), nil)
			if err != nil {
				t.Fatal(err)
			}

			if bytes.Equal(buf.Bytes(), out) {
				return
			}

			t.Error("not correctly formatted")

			if diff, err := testDiff(buf.String(), string(out)); err != nil {
				t.Error(err)
			} else {
				t.Log(diff)
			}
		})
	}
}
