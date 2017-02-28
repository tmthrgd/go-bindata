// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/tmthrgd/go-bindata"
)

func must(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "go-bindata: %v\n", err)
	os.Exit(1)
}

func main() {
	c, output := parseArgs()

	g, err := bindata.New(c)
	must(err)

	for i := 0; i < flag.NArg(); i++ {
		must(g.FindFiles(parseInput(flag.Arg(i))))
	}

	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	must(err)

	defer f.Close()

	_, err = g.WriteTo(f)
	must(err)
}

// parseArgs create s a new, filled configuration instance
// by reading and parsing command line options.
//
// This function exits the program with an error, if
// any of the command line options are incorrect.
func parseArgs() (c *bindata.Config, output string) {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <input directories>\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	var version bool
	flag.BoolVar(&version, "version", false, "Displays version information.")

	flag.StringVar(&output, "o", "./bindata.go", "Optional name of the output file to be generated.")

	c = &bindata.Config{
		Package:        "main",
		MemCopy:        true,
		Compress:       true,
		Metadata:       true,
		Restore:        true,
		HashLength:     16,
		AssetDir:       true,
		DecompressOnce: true,
	}

	var mode uint
	flag.BoolVar(&c.Debug, "debug", c.Debug, "Do not embed the assets, but provide the embedding API. Contents will still be loaded from disk.")
	flag.BoolVar(&c.Dev, "dev", c.Dev, "Similar to debug, but does not emit absolute paths. Expects a rootDir variable to already exist in the generated code's package.")
	flag.StringVar(&c.Tags, "tags", c.Tags, "Optional set of build tags to include.")
	flag.StringVar(&c.Prefix, "prefix", c.Prefix, "Optional path prefix to strip off asset names.")
	flag.StringVar(&c.Package, "pkg", c.Package, "Package name to use in the generated code.")
	flag.BoolVar(&c.MemCopy, "memcopy", c.MemCopy, "Do not use a .rodata hack to get rid of unnecessary memcopies. Refer to the documentation to see what implications this carries.")
	flag.BoolVar(&c.Compress, "compress", c.Compress, "Assets will be GZIP compressed when this flag is specified.")
	flag.BoolVar(&c.Metadata, "metadata", c.Metadata, "Assets will preserve size, mode, and modtime info.")
	flag.UintVar(&mode, "mode", uint(c.Mode), "Optional file mode override for all files.")
	flag.Int64Var(&c.ModTime, "modtime", c.ModTime, "Optional modification unix timestamp override for all files.")
	flag.BoolVar(&c.Restore, "restore", c.Restore, "[Deprecated]: use github.com/tmthrgd/go-bindata/restore.")
	flag.Var((*hashFormatValue)(&c.HashFormat), "hash", "Optional the format of name hashing to apply.")
	flag.IntVar(&c.HashLength, "hashlen", c.HashLength, "Optional length of hashes to be generated.")
	flag.Var((*hashEncodingValue)(&c.HashEncoding), "hashenc", `Optional the encoding of the hash to use. (default "hex")`)
	flag.Var((*hexEncodingValue)(&c.HashKey), "hashkey", "Optional hexadecimal key to use to turn the BLAKE2B hashing into a MAC.")
	flag.BoolVar(&c.AssetDir, "assetdir", c.AssetDir, "Provide the AssetDir APIs.")
	flag.Var((*appendRegexValue)(&c.Ignore), "ignore", "Regex pattern to ignore")
	flag.BoolVar(&c.DecompressOnce, "once", c.DecompressOnce, "Only GZIP decompress the resource once.")

	// Deprecated options
	var noMemCopy, noCompress, noMetadata bool
	flag.BoolVar(&noMemCopy, "nomemcopy", !c.MemCopy, "[Deprecated]: use -memcpy=false.")
	flag.BoolVar(&noCompress, "nocompress", !c.Compress, "[Deprecated]: use -compress=false.")
	flag.BoolVar(&noMetadata, "nometadata", !c.Metadata, "[Deprecated]: use -metadata=false.")

	flag.Parse()

	if version {
		fmt.Fprintf(os.Stderr, "go-bindata (Go runtime %s).\n", runtime.Version())
		io.WriteString(os.Stderr, "Copyright (c) 2010-2013, Jim Teeuwen.\n")
		io.WriteString(os.Stderr, "Copyright (c) 2017, Tom Thorogood.\n")
		os.Exit(0)
	}

	// Make sure we have input paths.
	if flag.NArg() == 0 {
		io.WriteString(os.Stderr, "Missing <input dir>\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if output == "" {
		cwd, err := os.Getwd()
		must(err)

		output = filepath.Join(cwd, "bindata.go")
	}

	c.Mode = os.FileMode(mode)

	var pkgSet, outputSet bool
	var memcopySet, nomemcopySet bool
	var compressSet, nocompressSet bool
	var metadataSet, nometadataSet bool
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "pkg":
			pkgSet = true
		case "o":
			outputSet = true
		case "memcopy":
			memcopySet = true
		case "nomemcopy":
			nomemcopySet = true
		case "compress":
			compressSet = true
		case "nocompress":
			nocompressSet = true
		case "metadata":
			metadataSet = true
		case "nometadata":
			nometadataSet = true
		}
	})

	// Change pkg to containing directory of output. If output flag is set and package flag is not.
	if outputSet && !pkgSet {
		pkg := identifier(filepath.Base(filepath.Dir(output)))
		if pkg != "" {
			c.Package = pkg
		}
	}

	if !memcopySet && nomemcopySet {
		c.MemCopy = !noMemCopy
	}

	if !compressSet && nocompressSet {
		c.Compress = !noCompress
	}

	if !metadataSet && nometadataSet {
		c.Metadata = !noMetadata
	}

	if !c.MemCopy && c.Compress {
		io.WriteString(os.Stderr, "The use of -memcopy=false with -compress is deprecated.\n")
	}

	must(validateOutput(output))
	return
}

func validateOutput(output string) error {
	stat, err := os.Lstat(output)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		// File does not exist. This is fine, just make
		// sure the directory it is to be in exists.
		if dir, _ := filepath.Split(output); dir != "" {
			if err = os.MkdirAll(dir, 0744); err != nil {
				return err
			}
		}
	} else if stat.IsDir() {
		return errors.New("output path is a directory")
	}

	return nil
}

// parseInput determines whether the given path has a recursive indicator and
// returns a new path with the recursive indicator chopped off if it does.
//
//  ex:
//      /path/to/foo/...    -> (/path/to/foo, true)
//      /path/to/bar        -> (/path/to/bar, false)
func parseInput(input string) (path string, recursive bool) {
	return filepath.Clean(strings.TrimSuffix(input, "/...")),
		strings.HasSuffix(input, "/...")
}

// identifier removes all characters from a string that are not valid in
// an identifier according to the Go Programming Language Specification.
//
// The logic in the switch statement aws taken from go/source package:
// https://github.com/golang/go/blob/a1a688fa0012f7ce3a37e9ac0070461fe8e3f28e/src/go/scanner/scanner.go#L257-#L271
func identifier(val string) string {
	return strings.Map(func(ch rune) rune {
		switch {
		case 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' ||
			ch >= utf8.RuneSelf && unicode.IsLetter(ch):
			return ch
		case '0' <= ch && ch <= '9' ||
			ch >= utf8.RuneSelf && unicode.IsDigit(ch):
			return ch
		default:
			return -1
		}
	}, val)
}
