// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import "text/template"

func init() {
	template.Must(baseTemplate.New("restore").Parse(`// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	path := filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)

	data, info, err := AssetAndInfo(name)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	if err = ioutil.WriteFile(path, data, info.Mode()); err != nil {
		return err
	}

	return os.Chtimes(path, info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}

	// Dir
	for _, child := range children {
		if err = RestoreAssets(dir, filepath.Join(name, child)); err != nil {
			return err
		}
	}

	return nil
}`))
}
