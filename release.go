// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"path"
	"text/template"
	"unicode/utf8"
)

func init() {
	template.Must(baseTemplate.New("release").Funcs(template.FuncMap{
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
	}).Parse(`{{if $.Config.Compress -}}
import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
{{- if $.Config.Restore}}
	"io/ioutil"
{{- end}}
	"os"
{{- if $.Config.Restore}}
	"path/filepath"
{{- end}}
	"strings"
	"time"
)

func bindataRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
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

{{else -}}
import (
{{- if $.Config.Restore}}
	"io/ioutil"
{{- end}}
	"os"
{{- if $.Config.Restore}}
	"path/filepath"
{{- end}}
{{- if not $.Config.MemCopy}}
	"reflect"
{{- end}}
	"strings"
	"time"
{{- if not $.Config.MemCopy}}
	"unsafe"
{{- end}}
)

{{if not $.Config.MemCopy -}}
func bindataRead(data string) []byte {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&data))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(data)
	bx.Cap = bx.Len
	return b
}

{{end -}}
{{end -}}

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

func (fi *bindataFileInfo) Name() string {
	return fi.name
}
func (fi *bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi *bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi *bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi *bindataFileInfo) IsDir() bool {
	return false
}
func (fi *bindataFileInfo) Sys() interface{} {
	return nil
}
{{- if ne $.Config.HashFormat 0}}
func (fi *bindataFileInfo) OriginalName() string {
	return fi.original
}
func (fi *bindataFileInfo) FileHash() string {
	return fi.hash
}
{{- end}}

{{- if ne $.Config.HashFormat 0}}

type FileInfo interface {
	os.FileInfo

	OriginalName() string
	FileHash() string
}
{{- end}}

{{range $.Assets -}}
{{$data := read .Path -}}

var _bindata_{{.Func}} = {{if $.Config.Compress -}}
	"" +
	{{gzip $data "\t" 28}}
{{- else if $.Config.MemCopy -}}
	[]byte(
	{{- if and (utf8Valid $data) (not (containsZero $data)) -}}
		` + "`{{printf \"%s\" (sanitize $data)}}`" + `
	{{- else -}}
		{{printf "%+q" $data}}
	{{- end -}}
	)
{{- else -}}
	bindataRead("" +
	{{wrap $data "\t" 28 -}}
	)
{{- end}}

var _bininfo_{{.Func}} = &bindataFileInfo{
	name: {{printf "%q" .Name}},

{{- if $.Config.Metadata}}
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
	hash: {{wrap .Hash "\t\t" 26}},
{{- end}}
}

func {{.Func}}() ([]byte, os.FileInfo, error) {
{{- if $.Config.Compress}}
	data, err := bindataRead(
		_bindata_{{.Func}},
		{{printf "%q" .Name}},
	)
	if err != nil {
		return nil, nil, err
	}

	return data, _bininfo_{{.Func}}, nil
{{- else}}
	return _bindata_{{.Func}}, _bininfo_{{.Func}}, nil
{{- end}}
}

{{end}}`))
}
