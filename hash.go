// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"encoding/base32"
	"encoding/base64"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/tmthrgd/go-hex"
	"golang.org/x/crypto/blake2b"
)

var base32Enc = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")

// hashFile applies name hashing with a given format,
// length and encoding. It returns the hashed name, the
// hash and any error that occurred. The hash is a BLAKE2B
// digest of the file contents.
func (asset *binAsset) hashFile(opts *GenerateOptions) error {
	if opts.HashFormat == NoHash {
		return nil
	}

	h, err := blake2b.New512(opts.HashKey)
	if err != nil {
		return err
	}

	rc, err := asset.Open()
	if err != nil {
		return err
	}

	buf := getSizedBuffer(rc)

	_, err = io.CopyBuffer(h, rc, buf.Bytes()[:buf.Cap()])

	rc.Close()
	bufPool.Put(buf)

	if err != nil {
		return err
	}

	asset.Hash = h.Sum(nil)

	if opts.HashFormat == NameUnchanged {
		return nil
	}

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

	if l > uint(len(enc)) {
		return errors.New("invalid HashLength: longer than generated hash")
	}

	dir, file := filepath.Split(asset.Name)
	ext := filepath.Ext(file)
	enc = enc[:l]

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
