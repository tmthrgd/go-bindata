// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"
)

type loggerWriter struct{ testing.TB }

func (w *loggerWriter) Write(p []byte) (n int, err error) {
	w.Log(string(bytes.TrimRight(p, "\r\n")))
	return len(p), nil
}

func TestIssue8(t *testing.T) {
	// This test case covers https://github.com/tmthrgd/go-bindata/issues/8.
	//
	// It generates a file with ~43,000 string concatenations which
	// triggers https://golang.org/issue/16394.

	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	dir, err := ioutil.TempDir("", "go-bindata-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	f, err := os.Create(path.Join(dir, "issue-8.go"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err = (Files{
		&testFile{
			path: "issue-8.bin",
			size: 1 << 20,
		},
	}).Generate(f, &GenerateOptions{
		Package: "main",
	}); err != nil {
		t.Fatal(err)
	}

	if _, err = io.WriteString(f, "\nfunc main() {}\n"); err != nil {
		t.Fatal(err)
	}

	if err = f.Close(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "build", ".")

	cmd.Dir = dir
	cmd.Env = os.Environ()

	cmd.Stderr = &loggerWriter{t}
	cmd.Stdout = cmd.Stderr

	if err = cmd.Run(); err != nil {
		t.Fatal(err)
	}
}
