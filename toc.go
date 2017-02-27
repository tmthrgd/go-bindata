// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"io"
	"text/template"
)

// writeTOC writes the table of contents file.
func writeTOC(w io.Writer, toc []binAsset, hashFormat HashFormat) error {
	return tocTemplate.Execute(w, struct {
		AssetName bool
		Assets    []binAsset
	}{hashFormat != NoHash && hashFormat != NameUnchanged, toc})
}

var tocTemplate = template.Must(template.New("toc").Parse(`// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	data, _, err := AssetAndInfo(name)
	return data, err
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
	_, info, err := AssetAndInfo(name)
	return info, err
}

// AssetAndInfo loads and returns the asset and asset info for the
// given name. It returns an error if the asset could not be found
// or could not be loaded.
func AssetAndInfo(name string) ([]byte, os.FileInfo, error) {
	if f, ok := _bindata[strings.Replace(name, "\\", "/", -1)]; ok {
		return f()
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
var _bindata = map[string]func() ([]byte, os.FileInfo, error){
{{range .Assets}}	{{printf "%q" .Name}}: {{.Func}},
{{end -}}
}

{{- if $.AssetName}}

// AssetName returns the hashed name associated with an asset of a
// given name.
func AssetName(name string) (string, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if hashedName, ok := _hashNames[canonicalName]; ok {
		return hashedName, nil
	}
	return "", &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

var _hashNames = map[string]string{
{{range .Assets}}	{{printf "%q" .OriginalName}}: {{printf "%q" .Name}},
{{end -}}
}
{{- end}}
`))
