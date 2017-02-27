// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

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
func hashFile(c *Config, path, name string) (newName string, hash []byte, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	h, err := blake2b.New512(c.HashKey)
	if err != nil {
		return
	}

	if _, err = io.Copy(h, f); err != nil {
		return
	}

	hash = h.Sum(nil)

	if c.HashFormat == NameUnchanged {
		newName = name
		return
	}

	var enc string
	switch c.HashEncoding {
	case HexHash:
		enc = hex.EncodeToString(hash)
	case Base32Hash:
		enc = strings.TrimSuffix(base32Enc.EncodeToString(hash), "=")
	case Base64Hash:
		enc = base64.RawURLEncoding.EncodeToString(hash)
	default:
		panic("unreachable")
	}

	dir, file := filepath.Split(name)
	ext := filepath.Ext(file)
	enc = enc[:c.HashLength]

	switch c.HashFormat {
	case DirHash:
		newName = filepath.Join(dir, enc, file)
	case NameHashSuffix:
		file = strings.TrimSuffix(file, ext)
		newName = filepath.Join(dir, file+"-"+enc+ext)
	case HashWithExt:
		newName = filepath.Join(dir, enc+ext)
	default:
		panic("unreachable")
	}

	return
}
