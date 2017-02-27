// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"

	"github.com/tmthrgd/go-bindata"
)

type hashEncodingValue bindata.HashEncoding

func (he *hashEncodingValue) String() string {
	if he == nil {
		return "hex"
	}

	switch bindata.HashEncoding(*he) {
	case bindata.HexHash:
		return "hex"
	case bindata.Base32Hash:
		return "base32"
	case bindata.Base64Hash:
		return "base64"
	default:
		panic("invalid HashFormat")
	}
}

func (he *hashEncodingValue) Set(value string) error {
	switch strings.ToLower(value) {
	case "hex":
		*he = hashEncodingValue(bindata.HexHash)
	case "base32":
		*he = hashEncodingValue(bindata.Base32Hash)
	case "base64":
		*he = hashEncodingValue(bindata.Base64Hash)
	default:
		return fmt.Errorf("invalid value %s, expected one of: hex, base32 or base64", value)
	}

	return nil
}
