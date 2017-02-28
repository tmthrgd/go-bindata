// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Implement sort.Interface for []os.FileInfo based on Name()
type byName []os.FileInfo

func (v byName) Len() int           { return len(v) }
func (v byName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v byName) Less(i, j int) bool { return v[i].Name() < v[j].Name() }

// findFiles recursively finds all the file paths in the given directory tree.
// They are added to the given map as keys. Values will be safe function names
// for each file, which will be used when generating the output code.
func (g *Generator) findFiles(dir, prefix string, recursive bool) error {
	dirpath := dir
	if len(prefix) > 0 {
		dirpath, _ = filepath.Abs(dirpath)
		prefix, _ = filepath.Abs(prefix)
		prefix = filepath.ToSlash(prefix)
	}

	fi, err := os.Stat(dirpath)
	if err != nil {
		return err
	}

	var list []os.FileInfo

	if !fi.IsDir() {
		dirpath = filepath.Dir(dirpath)
		list = []os.FileInfo{fi}
	} else {
		g.visited[dirpath] = struct{}{}

		fd, err := os.Open(dirpath)
		if err != nil {
			return err
		}
		defer fd.Close()

		if list, err = fd.Readdir(0); err != nil {
			return err
		}

		// Sort to make output stable between invocations
		sort.Sort(byName(list))
	}

outer:
	for _, file := range list {
		var asset binAsset
		asset.Path = filepath.Join(dirpath, file.Name())
		asset.Name = filepath.ToSlash(asset.Path)

		for _, re := range g.c.Ignore {
			if re.MatchString(asset.Path) {
				continue outer
			}
		}

		if file.IsDir() {
			if !recursive {
				continue
			}

			g.visited[asset.Path] = struct{}{}

			path := filepath.Join(dir, file.Name())
			if err = g.findFiles(path, prefix, recursive); err != nil {
				return err
			}

			continue
		} else if file.Mode()&os.ModeSymlink == os.ModeSymlink {
			linkPath, err := os.Readlink(asset.Path)
			if err != nil {
				return err
			}

			if !filepath.IsAbs(linkPath) {
				if linkPath, err = filepath.Abs(dirpath + "/" + linkPath); err != nil {
					return err
				}
			}

			if _, ok := g.visited[linkPath]; ok {
				continue
			}

			g.visited[linkPath] = struct{}{}

			if err = g.findFiles(asset.Path, prefix, recursive); err != nil {
				return err
			}

			continue
		}

		if strings.HasPrefix(asset.Name, prefix) {
			asset.Name = asset.Name[len(prefix):]
		} else if strings.HasSuffix(dir, file.Name()) {
			// Issue 110: dir is a full path, including
			// the file name (minus the basedir), so this
			// is what we have to use.
			asset.Name = dir
		} else {
			// Issue 110: dir is just that, a plain
			// directory, so we have to add the file's
			// name to it to form the full asset path.
			asset.Name = filepath.Join(dir, file.Name())
		}

		// If we have a leading slash, get rid of it.
		asset.Name = strings.TrimPrefix(asset.Name, "/")

		// This shouldn't happen.
		if len(asset.Name) == 0 {
			return fmt.Errorf("Invalid file: %v", asset.Path)
		}

		if g.c.HashFormat != NoHash {
			asset.OriginalName = asset.Name

			if err = hashFile(&g.c, &asset); err != nil {
				return err
			}
		}

		asset.Path, _ = filepath.Abs(asset.Path)
		g.toc = append(g.toc, asset)
	}

	return nil
}
