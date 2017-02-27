// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import "text/template"

func init() {
	template.Must(baseTemplate.New("debug").Parse(`import (
	"io/ioutil"
	"os"
{{- if or $.Config.Dev $.Config.Restore}}
	"path/filepath"
{{- end}}
	"strings"
)

// AssetAndInfo loads and returns the asset and asset info for the
// given name. It returns an error if the asset could not be found
// or could not be loaded.
func AssetAndInfo(name string) ([]byte, os.FileInfo, error) {
	path, ok := _bindata[strings.Replace(name, "\\", "/", -1)]
	if !ok {
		return nil, nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}

{{- if $.Config.Dev}}

	path = filepath.Join(rootDir, path)
{{- end}}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		return nil, nil, err
	}

	return data, fi, nil
}

// _bindata is a table, mapping each file to its path.
var _bindata = map[string]string{
{{$max := maxNameLength .Assets -}}
{{range .Assets}}	{{printf "%q" .Name}}:
	{{- repeat " " (sub $max (len .Name))}} {{if $.Config.Dev -}}
	{{printf "%q" .Name}}
{{- else -}}
	{{printf "%q" .Path}}
{{- end}},
{{end -}}
}`))
}
