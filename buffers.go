// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"sync"
)

var bufPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func getSizedBuffer(f file) *bytes.Buffer {
	var n int
	if fi, err := f.Stat(); err == nil {
		// Don't preallocate a huge buffer, just in case.
		if size := fi.Size(); size < 1e9 {
			n = int(size)
		}
	}

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Grow(n + bytes.MinRead)
	return buf
}
