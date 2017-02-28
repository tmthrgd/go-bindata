// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"io"
	"strings"
	"text/template"
)

// Translate reads assets from an input directory, converts
// them to Go code and writes the generated code to w.
func Translate(w io.Writer, c *Config) (err error) {
	// Ensure our configuration has sane values.
	if err = c.validate(); err != nil {
		return
	}

	var toc []binAsset
	var visitedPaths = make(map[string]struct{})

	// Locate all the assets.
	for _, input := range c.Input {
		if err = findFiles(c, input.Path, c.Prefix, input.Recursive, &toc, visitedPaths); err != nil {
			return
		}
	}

	return baseTemplate.Execute(w, struct {
		Config    *Config
		AssetName bool
		Assets    []binAsset
	}{c, c.HashFormat != NoHash && c.HashFormat != NameUnchanged, toc})
}

var baseTemplate = template.Must(template.New("base").Funcs(template.FuncMap{
	"repeat": strings.Repeat,
	"sub": func(a, b int) int {
		return a - b
	},
}).Parse(`
{{- template "header" .}}

{{if or $.Config.Debug $.Config.Dev -}}
{{- template "debug" . -}}
{{- else -}}
{{- template "release" . -}}
{{- end}}

{{template "common" . -}}

{{- if $.Config.AssetDir}}

{{template "tree" . -}}
{{- end}}
`))
