// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"encoding/base32"
	"encoding/base64"
	"errors"
	"hash"
	"io"
	"path/filepath"
	"strings"

	"github.com/tmthrgd/go-hex"
)

var base32Enc = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")

// hashFile hashes the asset and returns the resulting hash.
func (asset *binAsset) hashFile(h hash.Hash) ([]byte, error) {
	rc, err := asset.Open()
	if err != nil {
		return nil, err
	}

	buf := getSizedBuffer(rc)

	_, err = io.CopyBuffer(h, rc, buf.Bytes()[:buf.Cap()])

	rc.Close()
	bufPool.Put(buf)

	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// mangleName applies name hashing with a given format,
// length and encoding. It replaces asset.Name with the
// mangled name.
func (asset *binAsset) mangleName(opts *GenerateOptions) error {
	var enc string
	switch opts.HashEncoding {
	case HexHash:
		enc = hex.EncodeToString(asset.Hash)
	case Base32Hash:
		enc = strings.TrimSuffix(base32Enc.EncodeToString(asset.Hash), "=")
	case Base64Hash:
		enc = base64.RawURLEncoding.EncodeToString(asset.Hash)
	default:
		return errors.New("invalid HashEncoding")
	}

	l := opts.HashLength
	if l == 0 {
		l = 16
	}

	if l < uint(len(enc)) {
		enc = enc[:l]
	}

	dir, file := filepath.Split(asset.Name)
	ext := filepath.Ext(file)

	switch opts.HashFormat {
	case DirHash:
		asset.Name = filepath.Join(dir, enc, file)
	case NameHashSuffix:
		file = strings.TrimSuffix(file, ext)
		asset.Name = filepath.Join(dir, file+"-"+enc+ext)
	case HashWithExt:
		asset.Name = filepath.Join(dir, enc+ext)
	default:
		return errors.New("invalid HashFormat")
	}

	return nil
}
