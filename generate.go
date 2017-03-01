// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"hash"
	"io"
	"text/template"

	"golang.org/x/crypto/blake2b"
)

// binAsset holds information about a single asset to be processed.
type binAsset struct {
	File

	opts        *GenerateOptions
	Hash        []byte // Generated hash of file.
	mangledName string
}

// Generate writes the generated Go code to w.
func (f Files) Generate(w io.Writer, opts *GenerateOptions) error {
	if opts == nil {
		opts = &GenerateOptions{Package: "main"}
	}

	err := opts.validate()
	if err != nil {
		return err
	}

	var h hash.Hash
	if opts.HashFormat != NoHash {
		if h, err = blake2b.New512(opts.HashKey); err != nil {
			return err
		}
	}

	assets := make([]binAsset, 0, len(f))
	for i, file := range f {
		asset := binAsset{
			File: file,

			opts: opts,
		}

		if opts.HashFormat != NoHash {
			if i != 0 {
				h.Reset()
			}

			if err = asset.copy(h); err != nil {
				return err
			}

			asset.Hash = h.Sum(nil)
		}

		assets = append(assets, asset)
	}

	return baseTemplate.Execute(w, struct {
		*GenerateOptions
		Assets []binAsset
	}{opts, assets})
}

var baseTemplate = template.Must(template.New("base").Parse(`
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
