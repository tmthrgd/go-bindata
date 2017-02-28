// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"path"
	"text/template"
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
		"maxOriginalNameLength": func(toc []binAsset) int {
			l := 0
			for _, asset := range toc {
				if len(asset.OriginalName) > l {
					l = len(asset.OriginalName)
				}
			}

			return l
		},
	}).Parse(`
{{- $unsafeRead := and (not $.Config.Compress) (not $.Config.MemCopy) -}}
import (
{{- if $.Config.Compress}}
	"bytes"
	"compress/gzip"
	"io"
{{- end}}
	"os"
{{- if $unsafeRead}}
	"reflect"
{{- end}}
	"strings"
{{- if and $.Config.Compress $.Config.DecompressOnce}}
	"sync"
{{- end}}
	"time"
{{- if $unsafeRead}}
	"unsafe"
{{- end}}
{{- if $.Config.Restore}}

	"github.com/tmthrgd/go-bindata/restore"
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

{{else if not $.Config.Compress -}}
type bindataRead []byte

{{end -}}

type asset struct {
	name string
{{- if $.AssetName}}
	orig string
{{- end -}}
{{- if $.Config.Compress}}
	data string
	size int64
{{- else}}
	data []byte
{{- end -}}
{{- if and $.Config.Metadata (le $.Config.Mode 0)}}
	mode os.FileMode
{{- end -}}
{{- if and $.Config.Metadata (le $.Config.ModTime 0)}}
	time time.Time
{{- end -}}
{{- if ne $.Config.HashFormat 0}}
	hash []byte
{{- end}}
{{- if and $.Config.Compress $.Config.DecompressOnce}}

	once  sync.Once
	bytes []byte
	err   error
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
	return a.time
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
	return a.orig
{{- else}}
	return a.name
{{- end}}
}

func (a *asset) FileHash() []byte {
	return a.hash
}

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
	{{- if $.AssetName}}
		orig: {{printf "%q" .OriginalName}},
	{{- end}}
		data: {{$data := read .Path -}}
		{{- if $.Config.Compress -}}
			"" +
			{{gzip $data "\t\t\t" 24}}
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

	{{- if and $.Config.Metadata (le $.Config.ModTime 0)}}
		{{$mod := (stat .Path).ModTime -}}
		time: time.Unix({{$mod.Unix}}, {{$mod.Nanosecond}}),
	{{- end -}}

	{{- if ne $.Config.HashFormat 0}}
	{{- if $.Config.Compress}}
		hash: []byte("" +
	{{- else}}
		hash: bindataRead("" +
	{{- end}}
			{{wrap .Hash "\t\t\t" 24 -}}
		),
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
{{if and $.Config.Compress $.Config.DecompressOnce}}
	a.once.Do(func() {
		var gz *gzip.Reader
		if gz, a.err = gzip.NewReader(strings.NewReader(a.data)); a.err != nil {
			return
		}

		var buf bytes.Buffer
		if _, a.err = io.Copy(&buf, gz); a.err != nil {
			return
		}

		if a.err = gz.Close(); a.err == nil {
			a.bytes = buf.Bytes()
		}
	})
	if a.err != nil {
		return nil, nil, &os.PathError{Op: "read", Path: name, Err: a.err}
	}

	return a.bytes, a, nil
{{- else if $.Config.Compress}}
	gz, err := gzip.NewReader(strings.NewReader(a.data))
	if err != nil {
		return nil, nil, &os.PathError{Op: "read", Path: name, Err: err}
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, gz); err != nil {
		return nil, nil, &os.PathError{Op: "read", Path: name, Err: err}
	}

	if err = gz.Close(); err != nil {
		return nil, nil, &os.PathError{Op: "read", Path: name, Err: err}
	}

	return buf.Bytes(), a, nil
{{- else}}
	return a.data, a, nil
{{- end}}
}

{{- if $.AssetName}}

var _hashNames = map[string]string{
{{$max := maxOriginalNameLength .Assets -}}
{{range .Assets}}	{{printf "%q" .OriginalName}}:
	{{- repeat " " (sub $max (len .OriginalName))}} {{printf "%q" .Name}},
{{end -}}
}

// AssetName returns the hashed name associated with an asset of a
// given name.
func AssetName(name string) (string, error) {
	if name, ok := _hashNames[strings.Replace(name, "\\", "/", -1)]; ok {
		return name, nil
	}

	return "", &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}
{{- end}}`))
}
