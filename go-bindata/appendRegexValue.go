// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"bytes"
	"regexp"
)

type appendRegexValue []*regexp.Regexp

func (ar *appendRegexValue) String() string {
	if ar == nil {
		return ""
	}

	var buf bytes.Buffer

	for i, r := range *ar {
		if i != 0 {
			buf.WriteString(", ")
		}

		buf.WriteString(r.String())
	}

	return buf.String()
}

func (ar *appendRegexValue) Set(value string) error {
	r, err := regexp.Compile(value)
	if err != nil {
		return err
	}

	if *ar == nil {
		*ar = make([]*regexp.Regexp, 0, 1)
	}

	*ar = append(*ar, r)
	return nil
}
