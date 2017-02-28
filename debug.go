// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import "text/template"

func init() {
	template.Must(baseTemplate.New("debug").Funcs(template.FuncMap{
		"maxNameLength": func(toc []binAsset) int {
			l := 0
			for _, asset := range toc {
				if len(asset.Name) > l {
					l = len(asset.Name)
				}
			}

			return l
		},
	}).Parse(`import (
	"io/ioutil"
	"os"
	"path/filepath"
{{- if $.Config.AssetDir}}
	"strings"
{{- end}}
{{- if $.Config.Restore}}

	"github.com/tmthrgd/go-bindata/restore"
{{- end}}
)

// AssetAndInfo loads and returns the asset and asset info for the
// given name. It returns an error if the asset could not be found
// or could not be loaded.
func AssetAndInfo(name string) ([]byte, os.FileInfo, error) {
	path, ok := _bindata[filepath.ToSlash(name)]
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
	{{printf "%q" .AbsPath}}
{{- end}},
{{end -}}
}`))
}
