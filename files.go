// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FindFilesOptions defines a set of options to use
// when searching for files.
type FindFilesOptions struct {
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

	// Recursive defines whether subdirectories of Path
	// should be recursively included in the conversion.
	Recursive bool

	// Ignores any filenames matching the regex pattern specified, e.g.
	// path/to/file.ext will ignore only that file, or \\.gitignore
	// will match any .gitignore file.
	//
	// This parameter can be provided multiple times.
	Ignore []*regexp.Regexp
}

// FindFiles adds all files inside a directory to the
// generated output.
func (g *Generator) FindFiles(path string, opts *FindFilesOptions) error {
	if opts == nil {
		opts = new(FindFilesOptions)
	}

	return filepath.Walk(path, func(assetPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if !opts.Recursive && assetPath != path {
				return filepath.SkipDir
			}

			return nil
		}

		for _, re := range opts.Ignore {
			if re.MatchString(assetPath) {
				return nil
			}
		}

		name := strings.TrimPrefix(filepath.ToSlash(
			strings.TrimPrefix(assetPath, opts.Prefix)), "/")
		if name == "" {
			panic("should be impossible")
		}

		asset := binAsset{
			Path:         assetPath,
			Name:         name,
			OriginalName: name,
		}

		if g.c.Debug {
			asset.AbsPath, _ = filepath.Abs(assetPath)
		}

		if g.c.HashFormat != NoHash {
			if err = hashFile(&g.c, &asset); err != nil {
				return err
			}
		}

		g.toc = append(g.toc, asset)
		return nil
	})
}
