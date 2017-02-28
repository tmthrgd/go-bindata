// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

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

	return bindata.HashFormat(*hf).String()
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
