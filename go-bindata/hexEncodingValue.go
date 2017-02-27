// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

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
