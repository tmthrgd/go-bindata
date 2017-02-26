// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"fmt"
	"io"
)

// writeDebug writes the debug code file.
func writeDebug(w io.Writer, c *Config, toc []Asset) error {
	err := writeDebugHeader(w)
	if err != nil {
		return err
	}

	for i := range toc {
		err = writeDebugAsset(w, c, &toc[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// writeDebugHeader writes output file headers.
// This targets debug builds.
func writeDebugHeader(w io.Writer) error {
	_, err := io.WriteString(w, `import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// bindataRead reads the given file from disk. It returns an error on failure.
func bindataRead(path, name string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

`)
	return err
}

// writeDebugAsset write a debug entry for the given asset.
// A debug entry is simply a function which reads the asset from
// the original file (e.g.: from disk).
func writeDebugAsset(w io.Writer, c *Config, asset *Asset) error {
	pathExpr := fmt.Sprintf("%q", asset.Path)
	if c.Dev {
		pathExpr = fmt.Sprintf("filepath.Join(rootDir, %q)", asset.Name)
	}

	_, err := fmt.Fprintf(w, `// %s reads file data from disk. It returns an error on failure.
func %s() (*asset, error) {
	path := %s
	name := %q
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

`, asset.Func, asset.Func, pathExpr, asset.Name)
	return err
}
