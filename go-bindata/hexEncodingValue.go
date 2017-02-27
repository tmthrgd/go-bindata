// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import "encoding/hex"

type hexEncodingValue []byte

func (he *hexEncodingValue) String() string {
	if he == nil {
		return ""
	}

	return hex.EncodeToString(*he)
}

func (he *hexEncodingValue) Set(value string) (err error) {
	*he, err = hex.DecodeString(value)
	return
}
