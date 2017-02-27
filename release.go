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

type asset struct {
	data string
	info *bindataFileInfo
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

type asset struct {
	data []byte
	info *bindataFileInfo
}

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
	hash     []byte
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
func (fi *bindataFileInfo) FileHash() []byte {
	return fi.hash
}
{{- end}}

{{- if ne $.Config.HashFormat 0}}

type FileInfo interface {
	os.FileInfo

	OriginalName() string
	FileHash() []byte
}
{{- end}}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]asset{
{{range $.Assets}}	{{printf "%q" .Name}}: {
		{{$data := read .Path -}}
		{{- if $.Config.Compress -}}
			"" +
			{{gzip $data "\t\t\t" 24}}
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
			{{wrap $data "\t\t\t" 24 -}}
			)
		{{- end}},
		&bindataFileInfo{
			name: {{printf "%q" .Name}},

	{{- if $.Config.Metadata -}}
		{{- $info := stat .Path}}

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
			hash: []byte("" +
				{{wrap .Hash "\t\t\t\t" 22 -}}
			),
	{{- end}}
		},
	},
{{end -}}
}

// AssetAndInfo loads and returns the asset and asset info for the
// given name. It returns an error if the asset could not be found
// or could not be loaded.
func AssetAndInfo(name string) ([]byte, os.FileInfo, error) {
	f, ok := _bindata[strings.Replace(name, "\\", "/", -1)]
	if !ok {
		return nil, nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}
{{- if $.Config.Compress}}

	gz, err := gzip.NewReader(strings.NewReader(f.data))
	if err != nil {
		return nil, nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, gz); err != nil {
		return nil, nil, fmt.Errorf("Read %q: %v", name, err)
	}

	if err = gz.Close(); err != nil {
		return nil, nil, err
	}

	return buf.Bytes(), f.info, nil
{{- else}}

	return f.data, f.info, nil
{{- end}}
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
{{- end}}`))
}
