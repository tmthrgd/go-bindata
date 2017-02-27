// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"io"
	"text/template"
)

// writeDebug writes the debug code file.
func writeDebug(w io.Writer, c *Config, toc []binAsset) error {
	return debugTemplate.Execute(w, struct {
		Config *Config
		Assets []binAsset
	}{c, toc})
}

var debugTemplate = template.Must(template.New("debug").Parse(`import (
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
`))
