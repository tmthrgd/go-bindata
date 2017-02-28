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
	"sort"
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
	genOpts, findOpts, output := parseArgs()

	var all bindata.Files

	for i := 0; i < flag.NArg(); i++ {
		var path string
		path, findOpts.Recursive = parseInput(flag.Arg(i))

		files, err := bindata.FindFiles(path, findOpts)
		must(err)

		all = append(all, files...)
	}

	sort.Sort(all)

	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	must(err)

	defer f.Close()

	must(all.Generate(f, genOpts))
}

// parseArgs create s a new, filled configuration instance
// by reading and parsing command line options.
//
// This function exits the program with an error, if
// any of the command line options are incorrect.
func parseArgs() (genOpts *bindata.GenerateOptions, findOpts *bindata.FindFilesOptions, output string) {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <input directories>\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	var version bool
	flag.BoolVar(&version, "version", false, "Displays version information.")

	flag.StringVar(&output, "o", "./bindata.go", "Optional name of the output file to be generated.")

	genOpts = &bindata.GenerateOptions{
		Package:        "main",
		MemCopy:        true,
		Compress:       true,
		Metadata:       true,
		Restore:        true,
		HashLength:     16,
		AssetDir:       true,
		DecompressOnce: true,
	}
	findOpts = new(bindata.FindFilesOptions)

	var mode uint
	flag.BoolVar(&genOpts.Debug, "debug", genOpts.Debug, "Do not embed the assets, but provide the embedding API. Contents will still be loaded from disk.")
	flag.BoolVar(&genOpts.Dev, "dev", genOpts.Dev, "Similar to debug, but does not emit absolute paths. Expects a rootDir variable to already exist in the generated code's package.")
	flag.StringVar(&genOpts.Tags, "tags", genOpts.Tags, "Optional set of build tags to include.")
	flag.StringVar(&findOpts.Prefix, "prefix", "", "Optional path prefix to strip off asset names.")
	flag.StringVar(&genOpts.Package, "pkg", genOpts.Package, "Package name to use in the generated code.")
	flag.BoolVar(&genOpts.MemCopy, "memcopy", genOpts.MemCopy, "Do not use a .rodata hack to get rid of unnecessary memcopies. Refer to the documentation to see what implications this carries.")
	flag.BoolVar(&genOpts.Compress, "compress", genOpts.Compress, "Assets will be GZIP compressed when this flag is specified.")
	flag.BoolVar(&genOpts.Metadata, "metadata", genOpts.Metadata, "Assets will preserve size, mode, and modtime info.")
	flag.UintVar(&mode, "mode", uint(genOpts.Mode), "Optional file mode override for all files.")
	flag.Int64Var(&genOpts.ModTime, "modtime", genOpts.ModTime, "Optional modification unix timestamp override for all files.")
	flag.BoolVar(&genOpts.Restore, "restore", genOpts.Restore, "[Deprecated]: use github.com/tmthrgd/go-bindata/restore.")
	flag.Var((*hashFormatValue)(&genOpts.HashFormat), "hash", "Optional the format of name hashing to apply.")
	flag.UintVar(&genOpts.HashLength, "hashlen", genOpts.HashLength, "Optional length of hashes to be generated.")
	flag.Var((*hashEncodingValue)(&genOpts.HashEncoding), "hashenc", `Optional the encoding of the hash to use. (default "hex")`)
	flag.Var((*hexEncodingValue)(&genOpts.HashKey), "hashkey", "Optional hexadecimal key to use to turn the BLAKE2B hashing into a MAC.")
	flag.BoolVar(&genOpts.AssetDir, "assetdir", genOpts.AssetDir, "Provide the AssetDir APIs.")
	flag.Var((*appendRegexValue)(&findOpts.Ignore), "ignore", "Regex pattern to ignore")
	flag.BoolVar(&genOpts.DecompressOnce, "once", genOpts.DecompressOnce, "Only GZIP decompress the resource once.")

	// Deprecated options
	var noMemCopy, noCompress, noMetadata bool
	flag.BoolVar(&noMemCopy, "nomemcopy", !genOpts.MemCopy, "[Deprecated]: use -memcpy=false.")
	flag.BoolVar(&noCompress, "nocompress", !genOpts.Compress, "[Deprecated]: use -compress=false.")
	flag.BoolVar(&noMetadata, "nometadata", !genOpts.Metadata, "[Deprecated]: use -metadata=false.")

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

	genOpts.Mode = os.FileMode(mode)

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
			genOpts.Package = pkg
		}
	}

	if !memcopySet && nomemcopySet {
		genOpts.MemCopy = !noMemCopy
	}

	if !compressSet && nocompressSet {
		genOpts.Compress = !noCompress
	}

	if !metadataSet && nometadataSet {
		genOpts.Metadata = !noMetadata
	}

	if !genOpts.MemCopy && genOpts.Compress {
		io.WriteString(os.Stderr, "The use of -memcopy=false with -compress is deprecated.\n")
	}

	must(validateOutput(output))
	return
}

func validateOutput(output string) error {
	stat, err := os.Lstat(output)
	if err == nil {
		if stat.IsDir() {
			return errors.New("output path is a directory")
		}

		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	// File does not exist. This is fine, just make
	// sure the directory it is to be in exists.
	if dir, _ := filepath.Split(output); dir != "" {
		return os.MkdirAll(dir, 0744)
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
