// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"strings"
	"text/template"
)

var printerConfig = printer.Config{
	Mode:     printer.UseSpaces | printer.TabIndent,
	Tabwidth: 8,
}

// binAsset holds information about a single asset to be processed.
type binAsset struct {
	Path         string // Relative path.
	Name         string // Key used in TOC -- name by which asset is referenced.
	OriginalName string // Original Name before hashing applied to Name.
	Hash         []byte // Generated hash of file.
}

// Generate writes the generated Go code to w.
func (f Files) Generate(w io.Writer, opts *GenerateOptions) error {
	if opts == nil {
		opts = &GenerateOptions{Package: "main"}
	}

	if err := opts.validate(); err != nil {
		return err
	}

	assets := make([]binAsset, 0, len(f))
	for _, file := range f {
		asset := binAsset{
			Path:         file.Path,
			Name:         file.Name,
			OriginalName: file.Name,
		}
		if err := asset.hashFile(opts); err != nil {
			return err
		}

		assets = append(assets, asset)
	}

	ow := w
	if opts.Format {
		buf := bufPool.Get().(*bytes.Buffer)
		defer func() {
			buf.Reset()
			bufPool.Put(buf)
		}()
		w = buf
	}

	if err := baseTemplate.Execute(w, struct {
		*GenerateOptions
		Assets []binAsset
	}{opts, assets}); err != nil {
		return err
	}

	if opts.Format {
		fset := token.NewFileSet()

		f, err := parser.ParseFile(fset, "", w, parser.ParseComments)
		if err != nil {
			return err
		}

		return printerConfig.Fprint(ow, fset, f)
	}

	return nil
}

var baseTemplate = template.Must(template.New("base").Funcs(template.FuncMap{
	"repeat": strings.Repeat,
	"sub": func(a, b int) int {
		return a - b
	},
}).Parse(`
{{- template "header" .}}

{{if or $.Debug $.Dev -}}
{{- template "debug" . -}}
{{- else -}}
{{- template "release" . -}}
{{- end}}

{{template "common" . -}}

{{- if $.AssetDir}}

{{template "tree" . -}}
{{- end}}
`))
