// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/blake2b"
)

var base32Enc = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")

// hashFile applies name hashing with a given format,
// length and encoding. It returns the hashed name, the
// hash and any error that occurred. The hash is a BLAKE2B
// digest of the file contents.
func hashFile(c *Config, asset *binAsset) error {
	h, err := blake2b.New512(c.HashKey)
	if err != nil {
		return err
	}

	f, err := os.Open(asset.Path)
	if err != nil {
		return err
	}

	_, err = io.Copy(h, f)
	f.Close()
	if err != nil {
		return err
	}

	asset.Hash = h.Sum(nil)

	if c.HashFormat == NameUnchanged {
		return nil
	}

	var enc string
	switch c.HashEncoding {
	case HexHash:
		enc = hex.EncodeToString(asset.Hash)
	case Base32Hash:
		enc = strings.TrimSuffix(base32Enc.EncodeToString(asset.Hash), "=")
	case Base64Hash:
		enc = base64.RawURLEncoding.EncodeToString(asset.Hash)
	default:
		panic("unreachable")
	}

	dir, file := filepath.Split(asset.Name)
	ext := filepath.Ext(file)
	enc = enc[:c.HashLength]

	switch c.HashFormat {
	case DirHash:
		asset.Name = filepath.Join(dir, enc, file)
	case NameHashSuffix:
		file = strings.TrimSuffix(file, ext)
		asset.Name = filepath.Join(dir, file+"-"+enc+ext)
	case HashWithExt:
		asset.Name = filepath.Join(dir, enc+ext)
	default:
		panic("unreachable")
	}

	return nil
}
