// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"text/template"
	"unicode/utf8"
)

// writeRelease writes the release code file.
func writeRelease(w io.Writer, c *Config, toc []binAsset) error {
	return releaseTemplate.Execute(w, struct {
		Config *Config
		Assets []binAsset
	}{c, toc})
}

var releaseTemplate = template.Must(template.New("debug").Funcs(template.FuncMap{
	"stat": os.Stat,
	"read": ioutil.ReadFile,
	"name": func(name string) string {
		_, name = path.Split(name)
		return name
	},
	"wrap": func(data []byte, indent string, wrapAt int) string {
		var buf bytes.Buffer
		buf.WriteString(`"`)

		sw := &stringWriter{
			Writer: &buf,
			Indent: indent,
			WrapAt: wrapAt,
		}
		sw.Write(data)

		buf.WriteString(`"`)
		return buf.String()
	},
	"gzip": func(data []byte, indent string, wrapAt int) (string, error) {
		var buf bytes.Buffer
		buf.WriteString(`"`)

		gz := gzip.NewWriter(&stringWriter{
			Writer: &buf,
			Indent: indent,
			WrapAt: wrapAt,
		})

		if _, err := gz.Write(data); err != nil {
			return "", err
		}

		if err := gz.Close(); err != nil {
			return "", err
		}

		buf.WriteString(`"`)
		return buf.String(), nil
	},
	// sanitize prepares a valid UTF-8 string as a raw string constant.
	// Based on https://code.google.com/p/go/source/browse/godoc/static/makestatic.go?repo=tools
	"sanitize": func(b []byte) []byte {
		// Replace ` with `+"`"+`
		b = bytes.Replace(b, []byte("`"), []byte("`+\"`\"+`"), -1)

		// Replace BOM with `+"\xEF\xBB\xBF"+`
		// (A BOM is valid UTF-8 but not permitted in Go source files.
		// I wouldn't bother handling this, but for some insane reason
		// jquery.js has a BOM somewhere in the middle.)
		return bytes.Replace(b, []byte("\xEF\xBB\xBF"), []byte("`+\"\\xEF\\xBB\\xBF\"+`"), -1)
	},
	"utf8Valid": utf8.Valid,
	"containsZero": func(data []byte) bool {
		return bytes.Contains(data, []byte{0})
	},
}).Parse(`{{if $.Config.NoCompress -}}
import (
	"os"
{{- if $.Config.NoMemCopy}}
	"reflect"
{{- end}}
	"strings"
	"time"
{{- if $.Config.NoMemCopy}}
	"unsafe"
{{- end}}
)

{{if $.Config.NoMemCopy -}}
func bindataRead(data, name string) ([]byte, error) {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&data))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(data)
	bx.Cap = bx.Len
	return b, nil
}

{{end -}}
{{else -}}
import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

{{- if $.Config.NoMemCopy}}

func bindataRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
{{else}}

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
{{end -}}
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, gz); err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	if err = gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

{{end -}}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

{{- if ne $.Config.HashFormat 0}}

type FileInfo interface {
	os.FileInfo

	OriginalName() string
	FileHash() string
}
{{- end}}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
{{if ne $.Config.HashFormat 0}}
	original string
	hash     string
{{end -}}
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}
{{- if ne $.Config.HashFormat 0}}
func (fi bindataFileInfo) OriginalName() string {
	return fi.original
}
func (fi bindataFileInfo) FileHash() string {
	return fi.hash
}
{{- end}}

{{range $.Assets -}}
{{$data := read .Path -}}

var _{{.Func}} = {{if and $.Config.NoMemCopy $.Config.NoCompress -}}
	"" +
	{{wrap $data "\t" 28}}
{{- else if $.Config.NoCompress -}}
	[]byte(
	{{- if and (utf8Valid $data) (not (containsZero $data)) -}}
		` + "`{{printf \"%s\" (sanitize $data)}}`" + `
	{{- else -}}
		{{printf "%+q" $data}}
	{{- end -}}
	)
{{- else if $.Config.NoMemCopy -}}
	"" +
	{{gzip $data "\t" 28}}
{{- else -}}
	[]byte("" +
	{{gzip $data "\t" 28 -}}
	)
{{- end}}

func {{.Func}}() (*asset, error) {
{{- if and $.Config.NoCompress (not $.Config.NoMemCopy)}}
	bytes := []byte(_{{.Func}})
{{- else}}
	bytes, err := bindataRead(
		_{{.Func}},
		{{printf "%q" .Name}},
	)
	if err != nil {
		return nil, err
	}
{{end}}
	return &asset{
		bytes: bytes,
		info: bindataFileInfo{
			name: {{printf "%q" .Name}},

{{- if not $.Config.NoMetadata}}
	{{$info := stat .Path}}
			size:    {{$info.Size}},

	{{- if gt $.Config.Mode 0}}
			mode:    {{printf "%04o" $.Config.Mode}},
	{{- else}}
			mode:    {{printf "%04o" $info.Mode}},
	{{- end -}}

	{{- if gt $.Config.ModTime 0}}
			modTime: time.Unix($.Config.ModTime, 0),
	{{- else -}}
		{{$mod := $info.ModTime}}
			modTime: time.Unix({{$mod.Unix}}, {{$mod.Nanosecond}}),
	{{- end}}
{{- end}}

{{- if ne $.Config.HashFormat 0}}

			original: {{printf "%q" .OriginalName}},
			hash: {{wrap .Hash "\t\t\t\t" 22}},
{{- end}}
		},
	}, nil
}

{{end}}`))
