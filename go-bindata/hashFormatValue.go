package main

import (
	"fmt"
	"strings"

	"github.com/tmthrgd/go-bindata"
)

type hashFormatValue bindata.HashFormat

func (hf *hashFormatValue) String() string {
	switch bindata.HashFormat(*hf) {
	case bindata.NoHash:
		return ""
	case bindata.DirHash:
		return "dir"
	case bindata.NameHashSuffix:
		return "namesuffix"
	case bindata.HashWithExt:
		return "hashext"
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
	default:
		return fmt.Errorf("invalid value %s, expected one of: none, dir, namesuffix or hashext", value)
	}

	return nil
}
