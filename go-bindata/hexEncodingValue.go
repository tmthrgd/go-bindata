package main

import "encoding/hex"

type hexEncodingValue []byte

func (he *hexEncodingValue) String() string {
	return hex.EncodeToString(*he)
}

func (he *hexEncodingValue) Set(value string) (err error) {
	*he, err = hex.DecodeString(value)
	return
}
