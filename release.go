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
	}).Parse(`
{{- $unsafeRead := and (not $.Config.Compress) (not $.Config.MemCopy) -}}
import (
{{- if $.Config.Compress}}
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
{{- end}}
{{- if $.Config.Restore}}
	"io/ioutil"
{{- end}}
	"os"
{{- if $.Config.Restore}}
	"path/filepath"
{{- end}}
{{- if $unsafeRead}}
	"reflect"
{{- end}}
	"strings"
	"time"
{{- if $unsafeRead}}
	"unsafe"
{{- end}}
)

{{if $unsafeRead -}}
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

type asset struct {
	name string
{{- if $.Config.Compress}}
	data string
	size int64
{{- else}}
	data []byte
{{- end -}}

{{- if and $.Config.Metadata (le $.Config.Mode 0)}}
	mode os.FileMode
{{- end -}}

{{- if ne $.Config.HashFormat 0}}
	hash []byte
{{- end -}}

{{- if and $.Config.Metadata (le $.Config.ModTime 0)}}

	modTime time.Time
{{- end -}}

{{- if $.AssetName}}

	original string
{{- end}}
}

func (a *asset) Name() string {
	return a.name
}
func (a *asset) Size() int64 {
{{- if $.Config.Compress}}
	return a.size
{{- else}}
	return int64(len(a.data))
{{- end}}
}
func (a *asset) Mode() os.FileMode {
{{- if gt $.Config.Mode 0}}
	return {{printf "%04o" $.Config.Mode}}
{{- else if $.Config.Metadata}}
	return a.mode
{{- else}}
	return 0
{{- end}}
}
func (a *asset) ModTime() time.Time {
{{- if gt $.Config.ModTime 0}}
	return time.Unix({{$.Config.ModTime}}, 0)
{{- else if $.Config.Metadata}}
	return a.modTime
{{- else}}
	return time.Time{}
{{- end}}
}
func (*asset) IsDir() bool {
	return false
}
func (*asset) Sys() interface{} {
	return nil
}
{{- if ne $.Config.HashFormat 0}}
func (a *asset) OriginalName() string {
{{- if $.AssetName}}
	return a.original
{{- else}}
	return a.name
{{- end}}
}
func (a *asset) FileHash() []byte {
	return a.hash
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
var _bindata = map[string]*asset{
{{range $.Assets}}	{{printf "%q" .Name}}: &asset{
		name: {{printf "%q" (name .Name)}},
		data: {{$data := read .Path -}}
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
	{{- if $.Config.Compress}}
		size: {{len $data}},
	{{- end -}}

	{{- if and $.Config.Metadata (le $.Config.Mode 0)}}
		mode: {{printf "%04o" (stat .Path).Mode}},
	{{- end -}}

	{{- if ne $.Config.HashFormat 0}}
		hash: []byte("" +
			{{wrap .Hash "\t\t\t" 24 -}}
		),
	{{- end -}}

	{{- if and $.Config.Metadata (le $.Config.ModTime 0)}}

		{{$mod := (stat .Path).ModTime -}}
		modTime: time.Unix({{$mod.Unix}}, {{$mod.Nanosecond}}),
	{{- end -}}

	{{- if $.AssetName}}

		original: {{printf "%q" .OriginalName}},
	{{- end}}
	},
{{end -}}
}

// AssetAndInfo loads and returns the asset and asset info for the
// given name. It returns an error if the asset could not be found
// or could not be loaded.
func AssetAndInfo(name string) ([]byte, os.FileInfo, error) {
	a, ok := _bindata[strings.Replace(name, "\\", "/", -1)]
	if !ok {
		return nil, nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}
{{- if $.Config.Compress}}

	gz, err := gzip.NewReader(strings.NewReader(a.data))
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

	return buf.Bytes(), a, nil
{{- else}}

	return a.data, a, nil
{{- end}}
}

{{- if $.AssetName}}

// AssetName returns the hashed name associated with an asset of a
// given name.
func AssetName(name string) (string, error) {
	if name, ok := _hashNames[strings.Replace(name, "\\", "/", -1)]; ok {
		return name, nil
	}

	return "", &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

var _hashNames = map[string]string{
{{range .Assets}}	{{printf "%q" .OriginalName}}: {{printf "%q" .Name}},
{{end -}}
}
{{- end}}`))
}
