// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"regexp"

	"golang.org/x/crypto/blake2b"
)

// HashFormat specifies which format to use when hashing names.
type HashFormat int

const (
	// NoHash disables name hashing.
	NoHash HashFormat = iota
	// DirHash formats names like path/to/hash/name.ext.
	DirHash
	// NameHashSuffix formats names like path/to/name-hash.ext.
	NameHashSuffix
	// HashWithExt formats names like path/to/hash.ext.
	HashWithExt
	// NameUnchanged generates the file hash but does not change
	// the asset name.
	NameUnchanged
)

func (hf HashFormat) String() string {
	switch hf {
	case NoHash:
		return ""
	case DirHash:
		return "dir"
	case NameHashSuffix:
		return "namesuffix"
	case HashWithExt:
		return "hashext"
	case NameUnchanged:
		return "unchanged"
	default:
		return "unknown"
	}
}

// HashEncoding specifies which encoding to use when hashing names.
type HashEncoding int

const (
	// HexHash uses hexadecimal encoding.
	HexHash HashEncoding = iota
	// Base32Hash uses unpadded, lowercase standard base32
	// encoding (see RFC 4648).
	Base32Hash
	// Base64Hash uses an unpadded URL-safe base64 encoding
	// defined in RFC 4648.
	Base64Hash
)

func (he HashEncoding) String() string {
	switch he {
	case HexHash:
		return "hex"
	case Base32Hash:
		return "base32"
	case Base64Hash:
		return "base64"
	default:
		return "unknown"
	}
}

// InputConfig defines options on a asset directory to be convert.
type InputConfig struct {
	// Path defines a directory containing asset files to be included
	// in the generated output.
	Path string

	// Recursive defines whether subdirectories of Path
	// should be recursively included in the conversion.
	Recursive bool
}

// Config defines a set of options for the asset conversion.
type Config struct {
	// Name of the package to use. Defaults to 'main'.
	Package string

	// Tags specify a set of optional build tags, which should be
	// included in the generated output. The tags are appended to a
	// `// +build` line in the beginning of the output file
	// and must follow the build tags syntax specified by the go tool.
	Tags string

	// Input defines the directory path, containing all asset files as
	// well as whether to recursively process assets in any sub directories.
	Input []InputConfig

	// Prefix defines a path prefix which should be stripped from all
	// file names when generating the keys in the table of contents.
	// For example, running without the `-prefix` flag, we get:
	//
	// 	$ go-bindata /path/to/templates
	// 	go_bindata["/path/to/templates/foo.html"] = _path_to_templates_foo_html
	//
	// Running with the `-prefix` flag, we get:
	//
	// 	$ go-bindata -prefix "/path/to/" /path/to/templates/foo.html
	// 	go_bindata["templates/foo.html"] = templates_foo_html
	Prefix string

	// MemCopy will alter the way the output file is generated.
	//
	// If false, it will employ a hack that allows us to read the file data directly
	// from the compiled program's `.rodata` section. This ensures that when we call
	// call our generated function, we omit unnecessary mem copies.
	//
	// The downside of this, is that it requires dependencies on the `reflect` and
	// `unsafe` packages. These may be restricted on platforms like AppEngine and
	// thus prevent you from using this mode.
	//
	// Another disadvantage is that the byte slice we create, is strictly read-only.
	// For most use-cases this is not a problem, but if you ever try to alter the
	// returned byte slice, a runtime panic is thrown. Use this mode only on target
	// platforms where memory constraints are an issue.
	//
	// The default behaviour is to use the old code generation method. This
	// prevents the two previously mentioned issues, but will employ at least one
	// extra memcopy and thus increase memory requirements.
	//
	// For instance, consider the following two examples:
	//
	// This would be the default mode, using an extra memcopy but gives a safe
	// implementation without dependencies on `reflect` and `unsafe`:
	//
	// 	func myfile() []byte {
	// 		return []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a}
	// 	}
	//
	// Here is the same functionality, but uses the `.rodata` hack.
	// The byte slice returned from this example can not be written to without
	// generating a runtime error.
	//
	// 	var _myfile = "\x89\x50\x4e\x47\x0d\x0a\x1a"
	//
	// 	func myfile() []byte {
	// 		var empty [0]byte
	// 		sx := (*reflect.StringHeader)(unsafe.Pointer(&_myfile))
	// 		b := empty[:]
	// 		bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	// 		bx.Data = sx.Data
	// 		bx.Len = len(_myfile)
	// 		bx.Cap = bx.Len
	// 		return b
	// 	}
	MemCopy bool

	// Compress means the assets are GZIP compressed before being turned into
	// Go code. The generated function will automatically unzip the file data
	// when called. Defaults to true.
	Compress bool

	// Perform a debug build. This generates an asset file, which
	// loads the asset contents directly from disk at their original
	// location, instead of embedding the contents in the code.
	//
	// This is mostly useful if you anticipate that the assets are
	// going to change during your development cycle. You will always
	// want your code to access the latest version of the asset.
	// Only in release mode, will the assets actually be embedded
	// in the code. The default behaviour is Release mode.
	Debug bool

	// Perform a dev build, which is nearly identical to the debug option. The
	// only difference is that instead of absolute file paths in generated code,
	// it expects a variable, `rootDir`, to be set in the generated code's
	// package (the author needs to do this manually), which it then prepends to
	// an asset's name to construct the file path on disk.
	//
	// This is mainly so you can push the generated code file to a shared
	// repository.
	Dev bool

	// When false, size, mode and modtime are not preserved from files
	Metadata bool
	// When nonzero, use this as mode for all files.
	Mode os.FileMode
	// When nonzero, use this as unix timestamp for all files.
	ModTime int64

	// Ignores any filenames matching the regex pattern specified, e.g.
	// path/to/file.ext will ignore only that file, or \\.gitignore
	// will match any .gitignore file.
	//
	// This parameter can be provided multiple times.
	Ignore []*regexp.Regexp

	// [Deprecated]: use github.com/tmthrgd/go-bindata/restore.
	Restore bool

	// Which of the given name hashing formats to use.
	HashFormat HashFormat
	// The length of the hash to use, defaults to 16 characters.
	HashLength int
	// The encoding to use to encode the name hash.
	HashEncoding HashEncoding
	// The key to use to turn the BLAKE2B hashing into a MAC. Must be between
	// zero and 64 bytes long.
	HashKey []byte

	// When true, the AssetDir API will be provided.
	AssetDir bool

	// When true, only gzip decompress the data on first use.
	DecompressOnce bool
}

// NewConfig returns a default configuration struct.
func NewConfig() *Config {
	c := new(Config)
	c.Package = "main"
	c.MemCopy = true
	c.Compress = true
	c.Metadata = true
	c.Restore = true
	c.HashLength = 16
	c.AssetDir = true
	c.DecompressOnce = true
	return c
}

// validate ensures the config has sane values.
// Part of which means checking if certain file/directory paths exist.
func (c *Config) validate() error {
	if c == nil {
		return errors.New("go-bindata: Config not provided")
	}

	if len(c.Package) == 0 {
		return errors.New("go-bindata: missing package name")
	}

	for _, input := range c.Input {
		if _, err := os.Lstat(input.Path); err != nil {
			return err
		}
	}

	if c.Mode&^os.ModePerm != 0 {
		return errors.New("go-bindata: invalid mode specified")
	}

	if (c.Debug || c.Dev) && c.HashFormat != NoHash {
		return errors.New("go-bindata: HashFormat is not compatible with Debug and Dev")
	}

	length := 0
	switch c.HashEncoding {
	case HexHash:
		length = hex.EncodedLen(blake2b.Size)
	case Base32Hash:
		length = base32Enc.EncodedLen(blake2b.Size)
	case Base64Hash:
		length = base64.RawStdEncoding.EncodedLen(blake2b.Size)
	}

	if (c.HashFormat != NoHash && c.HashFormat != NameUnchanged) &&
		(c.HashLength <= 0 || c.HashLength > length) {
		return fmt.Errorf("go-bindata: HashLength must be between 1 and %d bytes in length", length)
	}

	if c.Restore && !c.AssetDir {
		return errors.New("go-bindata: Restore cannot be used without AssetDir")
	}

	return nil
}
