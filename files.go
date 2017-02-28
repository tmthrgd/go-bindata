// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"os"
	"path/filepath"
	"strings"
)

// FindFiles adds all files inside a directory to the
// generated output. If recursive is true, files within
// subdirectories of path will also be included.
func (g *Generator) FindFiles(path string, recursive bool) error {
	return filepath.Walk(path, func(assetPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if !recursive && assetPath != path {
				return filepath.SkipDir
			}

			return nil
		}

		for _, re := range g.c.Ignore {
			if re.MatchString(assetPath) {
				return nil
			}
		}

		name := strings.TrimPrefix(filepath.ToSlash(
			strings.TrimPrefix(assetPath, g.c.Prefix)), "/")
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
