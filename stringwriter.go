// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"fmt"
	"io"
)

type stringWriter struct {
	io.Writer
	Indent string
	WrapAt int
	c      int
}

func (w *stringWriter) Write(p []byte) (n int, err error) {
	var buf [4]byte
	buf[0] = '\\'
	buf[1] = 'x'

	for _, b := range p {
		const lowerHex = "0123456789abcdef"
		buf[2] = lowerHex[b/16]
		buf[3] = lowerHex[b%16]

		if _, err = w.Writer.Write(buf[:]); err != nil {
			return
		}

		n += 4
		w.c++

		if w.WrapAt == 0 || w.c%w.WrapAt != 0 {
			continue
		}

		nn, err := fmt.Fprintf(w.Writer, "\" +\n%s\"", w.Indent)
		if err != nil {
			return n, err
		}

		n += nn
	}

	return
}
