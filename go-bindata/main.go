// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tmthrgd/go-bindata"
)

func main() {
	if err := bindata.Translate(parseArgs()); err != nil {
		fmt.Fprintf(os.Stderr, "bindata: %v\n", err)
		os.Exit(1)
	}
}

// parseArgs create s a new, filled configuration instance
// by reading and parsing command line options.
//
// This function exits the program with an error, if
// any of the command line options are incorrect.
func parseArgs() *bindata.Config {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <input directories>\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	var version bool
	flag.BoolVar(&version, "version", false, "Displays version information.")

	c := bindata.NewConfig()
	flag.BoolVar(&c.Debug, "debug", c.Debug, "Do not embed the assets, but provide the embedding API. Contents will still be loaded from disk.")
	flag.BoolVar(&c.Dev, "dev", c.Dev, "Similar to debug, but does not emit absolute paths. Expects a rootDir variable to already exist in the generated code's package.")
	flag.StringVar(&c.Tags, "tags", c.Tags, "Optional set of build tags to include.")
	flag.StringVar(&c.Prefix, "prefix", c.Prefix, "Optional path prefix to strip off asset names.")
	flag.StringVar(&c.Package, "pkg", c.Package, "Package name to use in the generated code.")
	flag.BoolVar(&c.MemCopy, "memcopy", c.MemCopy, "Do not use a .rodata hack to get rid of unnecessary memcopies. Refer to the documentation to see what implications this carries.")
	flag.BoolVar(&c.Compress, "compress", c.Compress, "Assets will be GZIP compressed when this flag is specified.")
	flag.BoolVar(&c.Metadata, "metadata", c.Metadata, "Assets will preserve size, mode, and modtime info.")
	flag.UintVar(&c.Mode, "mode", c.Mode, "Optional file mode override for all files.")
	flag.Int64Var(&c.ModTime, "modtime", c.ModTime, "Optional modification unix timestamp override for all files.")
	flag.BoolVar(&c.Restore, "restore", c.Restore, "Provide the restore APIs.")
	flag.StringVar(&c.Output, "o", c.Output, "Optional name of the output file to be generated.")
	flag.Var((*hashFormatValue)(&c.HashFormat), "hash", "Optional the format of name hashing to apply.")
	flag.IntVar(&c.HashLength, "hashlen", c.HashLength, "Optional length of hashes to be generated.")
	flag.Var((*hashEncodingValue)(&c.HashEncoding), "hashenc", `Optional the encoding of the hash to use. (default "hex")`)
	flag.Var((*hexEncodingValue)(&c.HashKey), "hashkey", "Optional hexadecimal key to use to turn the BLAKE2B hashing into a MAC.")
	flag.BoolVar(&c.AssetDir, "assetdir", c.AssetDir, "Provide the AssetDir APIs.")
	flag.Var((*appendRegexValue)(&c.Ignore), "ignore", "Regex pattern to ignore")

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

	// Create input configurations.
	c.Input = make([]bindata.InputConfig, flag.NArg())
	for i := range c.Input {
		c.Input[i] = parseInput(flag.Arg(i))
	}

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
		if pkg := filepath.Base(filepath.Dir(c.Output)); pkg != "." && pkg != "/" {
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

	return c
}

// parseRecursive determines whether the given path has a recursive indicator and
// returns a new path with the recursive indicator chopped off if it does.
//
//  ex:
//      /path/to/foo/...    -> (/path/to/foo, true)
//      /path/to/bar        -> (/path/to/bar, false)
func parseInput(path string) bindata.InputConfig {
	return bindata.InputConfig{
		Path:      filepath.Clean(strings.TrimSuffix(path, "/...")),
		Recursive: strings.HasSuffix(path, "/..."),
	}
}
