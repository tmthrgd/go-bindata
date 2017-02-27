// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

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
		if err := writeTOCTree(&buf, toc); err != nil {
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

// Implement sort.Interface for []os.FileInfo based on Name()
type byName []os.FileInfo

func (v byName) Len() int           { return len(v) }
func (v byName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v byName) Less(i, j int) bool { return v[i].Name() < v[j].Name() }

// findFiles recursively finds all the file paths in the given directory tree.
// They are added to the given map as keys. Values will be safe function names
// for each file, which will be used when generating the output code.
func findFiles(c *Config, dir, prefix string, recursive bool, toc *[]binAsset, knownFuncs map[string]int, visitedPaths map[string]bool) error {
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
		visitedPaths[dirpath] = true
		fd, err := os.Open(dirpath)
		if err != nil {
			return err
		}

		defer fd.Close()

		list, err = fd.Readdir(0)
		if err != nil {
			return err
		}

		// Sort to make output stable between invocations
		sort.Sort(byName(list))
	}

	for _, file := range list {
		var asset binAsset
		asset.Path = filepath.Join(dirpath, file.Name())
		asset.Name = filepath.ToSlash(asset.Path)

		ignoring := false
		for _, re := range c.Ignore {
			if re.MatchString(asset.Path) {
				ignoring = true
				break
			}
		}
		if ignoring {
			continue
		}

		if file.IsDir() {
			if recursive {
				recursivePath := filepath.Join(dir, file.Name())
				visitedPaths[asset.Path] = true
				findFiles(c, recursivePath, prefix, recursive, toc, knownFuncs, visitedPaths)
			}
			continue
		} else if file.Mode()&os.ModeSymlink == os.ModeSymlink {
			var linkPath string
			if linkPath, err = os.Readlink(asset.Path); err != nil {
				return err
			}
			if !filepath.IsAbs(linkPath) {
				if linkPath, err = filepath.Abs(dirpath + "/" + linkPath); err != nil {
					return err
				}
			}
			if _, ok := visitedPaths[linkPath]; !ok {
				visitedPaths[linkPath] = true
				findFiles(c, asset.Path, prefix, recursive, toc, knownFuncs, visitedPaths)
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
		if len(asset.Name) > 0 && asset.Name[0] == '/' {
			asset.Name = asset.Name[1:]
		}

		// This shouldn't happen.
		if len(asset.Name) == 0 {
			return fmt.Errorf("Invalid file: %v", asset.Path)
		}

		if c.HashFormat != NoHash {
			asset.OriginalName = asset.Name
			asset.Name, asset.Hash, err = hashFile(c, asset.Path, asset.Name)
			if err != nil {
				return err
			}
		}

		asset.Func = safeFunctionName(asset.Name, knownFuncs)
		asset.Path, _ = filepath.Abs(asset.Path)
		*toc = append(*toc, asset)
	}

	return nil
}
