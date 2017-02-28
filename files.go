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

// File represents a single asset file.
type File struct {
	name string // Key used in TOC -- name by which asset is referenced.
	path string // Relative path.
	abs  string // Absolute path, only used for Debug.
}

// Files represents a collection of asset files.
type Files []*File

// Len implements sort.Interface.
func (f Files) Len() int {
	return len(f)
}

// Less implements sort.Interface.
func (f Files) Less(i, j int) bool {
	return f[i].name < f[j].name
}

// Swap implements sort.Interface.
func (f Files) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

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
func FindFiles(path string, opts *FindFilesOptions) (files Files, err error) {
	if opts == nil {
		opts = new(FindFilesOptions)
	}

	if err = filepath.Walk(path, func(assetPath string, info os.FileInfo, err error) error {
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

		abs, _ := filepath.Abs(assetPath)
		files = append(files, &File{
			name: name,
			path: assetPath,
			abs:  abs,
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return
}
