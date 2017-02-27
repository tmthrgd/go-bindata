// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"fmt"
	"strings"

	"github.com/tmthrgd/go-bindata"
)

type hashFormatValue bindata.HashFormat

func (hf *hashFormatValue) String() string {
	if hf == nil {
		return ""
	}

	switch bindata.HashFormat(*hf) {
	case bindata.NoHash:
		return ""
	case bindata.DirHash:
		return "dir"
	case bindata.NameHashSuffix:
		return "namesuffix"
	case bindata.HashWithExt:
		return "hashext"
	case bindata.NameUnchanged:
		return "unchanged"
	default:
		panic("invalid HashFormat")
	}
}

func (hf *hashFormatValue) Set(value string) error {
	switch strings.ToLower(value) {
	case "", "none":
		*hf = hashFormatValue(bindata.NoHash)
	case "dir":
		*hf = hashFormatValue(bindata.DirHash)
	case "namesuffix":
		*hf = hashFormatValue(bindata.NameHashSuffix)
	case "hashext":
		*hf = hashFormatValue(bindata.HashWithExt)
	case "unchanged":
		*hf = hashFormatValue(bindata.NameUnchanged)
	default:
		return fmt.Errorf("invalid value %s, expected one of: none, dir, namesuffix, hashext or unchanged", value)
	}

	return nil
}
