// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/tools/imports"
)

// Translate reads assets from an input directory, converts them
// to Go code and writes new files to the output specified
// in the given configuration.
func Translate(c *Config) error {
	var toc []binAsset

	// Ensure our configuration has sane values.
	err := c.validate()
	if err != nil {
		return err
	}

	var knownFuncs = make(map[string]int)
	var visitedPaths = make(map[string]bool)
	// Locate all the assets.
	for _, input := range c.Input {
		err = findFiles(c, input.Path, c.Prefix, input.Recursive, &toc, knownFuncs, visitedPaths)
		if err != nil {
			return err
		}
	}

	var buf bytes.Buffer

	// Write file header.
	if err := writeHeader(&buf, c, toc); err != nil {
		return err
	}

	// Write assets.
	if c.Debug || c.Dev {
		err = writeDebug(&buf, c, toc)
	} else {
		err = writeRelease(&buf, c, toc)
	}

	if err != nil {
		return err
	}

	// Write table of contents
	if err := writeTOC(&buf, toc, c.HashFormat); err != nil {
		return err
	}

	if c.AssetDir {
		// Write hierarchical tree of assets
		if err := writeTree(&buf, toc); err != nil {
			return err
		}
	}

	if c.Restore {
		// Write restore procedure
		if err := writeRestore(&buf); err != nil {
			return err
		}
	}

	out := buf.Bytes()
	if c.Format {
		if out, err = imports.Process(c.Output, out, nil); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(c.Output, out, 0666)
}
