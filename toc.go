// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"fmt"
	"io"
)

// writeTOC writes the table of contents file.
func writeTOC(w io.Writer, toc []binAsset, hashFormat HashFormat) error {
	err := writeTOCHeader(w)
	if err != nil {
		return err
	}

	for i := range toc {
		err = writeTOCAsset(w, &toc[i])
		if err != nil {
			return err
		}
	}

	if err := writeTOCFooter(w); err != nil {
		return err
	}

	if hashFormat == NoHash || hashFormat == NameUnchanged {
		return nil
	}

	if err := writeTOCHashNameHeader(w); err != nil {
		return err
	}

	for i := range toc {
		err = writeTOCHashNameAsset(w, &toc[i])
		if err != nil {
			return err
		}
	}

	return writeTOCHashNameFooter(w)
}

// writeTOCHeader writes the table of contents file header.
func writeTOCHeader(w io.Writer) error {
	_, err := io.WriteString(w, `// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, err
		}
		return a.bytes, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, err
		}
		return a.info, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// AssetAndInfo loads and returns the asset and asset info for the
// given name. It returns an error if the asset could not be found
// or could not be loaded.
func AssetAndInfo(name string) ([]byte, os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, nil, err
		}
		return a.bytes, a.info, nil
	}
	return nil, nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
`)
	return err
}

// writeTOCAsset write a TOC entry for the given asset.
func writeTOCAsset(w io.Writer, asset *binAsset) error {
	_, err := fmt.Fprintf(w, "\t%q: %s,\n", asset.Name, asset.Func)
	return err
}

// writeTOCFooter writes the table of contents file footer.
func writeTOCFooter(w io.Writer) error {
	_, err := io.WriteString(w, `}

`)
	return err
}

// writeTOCHashNameHeader writes the table of contents header for hash names.
func writeTOCHashNameHeader(w io.Writer) error {
	_, err := io.WriteString(w, `// AssetName returns the hashed name associated with an asset of a
// given name.
func AssetName(name string) (string, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if hashedName, ok := _hashNames[canonicalName]; ok {
		return hashedName, nil
	}
	return "", &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

var _hashNames = map[string]string{
`)
	return err
}

// writeTOCHashNameAsset write a hash name entry for the given asset.
func writeTOCHashNameAsset(w io.Writer, asset *binAsset) error {
	_, err := fmt.Fprintf(w, "\t%q: %q,\n", asset.OriginalName, asset.Name)
	return err
}

// writeTOCHashNameFooter writes the hash table of contents file footer.
func writeTOCHashNameFooter(w io.Writer) error {
	_, err := io.WriteString(w, `}

`)
	return err
}
