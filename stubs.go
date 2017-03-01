// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"io"
	"os"
	"path/filepath"
)

// for testing: path/filepath
var (
	abs  = filepath.Abs
	walk = filepath.Walk
)

// for testing: os
var (
	open = osOpen
	stat = os.Stat
)

type file interface {
	io.ReadCloser
	Stat() (os.FileInfo, error)
}

func osOpen(path string) (file, error) {
	// This is ok for our use, but beware of:
	// http://spf13.com/post/when-nil-is-not-nil/
	return os.Open(path)
}
