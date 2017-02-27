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

{{- range $.Assets}}

// {{.Func}} reads file data from disk. It returns an error on failure.
func {{.Func}}() ([]byte, os.FileInfo, error) {
{{- if $.Config.Dev}}
	path := filepath.Join(rootDir, {{printf "%q" .Name}})
{{- else}}
	path := {{printf "%q" .Path}}
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
{{- end}}

// AssetAndInfo loads and returns the asset and asset info for the
// given name. It returns an error if the asset could not be found
// or could not be loaded.
func AssetAndInfo(name string) ([]byte, os.FileInfo, error) {
	if f, ok := _bindata[strings.Replace(name, "\\", "/", -1)]; ok {
		return f()
	}

	return nil, nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, os.FileInfo, error){
{{$max := maxNameLength .Assets -}}
{{range .Assets}}	{{printf "%q" .Name}}:
	{{- repeat " " (sub $max (len .Name))}} {{.Func}},
{{end -}}
}`))
}
