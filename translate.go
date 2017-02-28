// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"io"
	"strings"
	"text/template"
)

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
		if err := hashFile(opts, &asset); err != nil {
			return err
		}

		assets = append(assets, asset)
	}

	return baseTemplate.Execute(w, struct {
		Opts      *GenerateOptions
		AssetName bool
		Assets    []binAsset
	}{opts, opts.HashFormat != NoHash && opts.HashFormat != NameUnchanged, assets})
}

var baseTemplate = template.Must(template.New("base").Funcs(template.FuncMap{
	"repeat": strings.Repeat,
	"sub": func(a, b int) int {
		return a - b
	},
}).Parse(`
{{- template "header" .}}

{{if or $.Opts.Debug $.Opts.Dev -}}
{{- template "debug" . -}}
{{- else -}}
{{- template "release" . -}}
{{- end}}

{{template "common" . -}}

{{- if $.Opts.AssetDir}}

{{template "tree" . -}}
{{- end}}
`))
