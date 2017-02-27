// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"os"
	"strings"
	"text/template"
)

// Translate reads assets from an input directory, converts them
// to Go code and writes new files to the output specified
// in the given configuration.
func Translate(c *Config) error {
	// Ensure our configuration has sane values.
	if err := c.validate(); err != nil {
		return err
	}

	var toc []binAsset
	var visitedPaths = make(map[string]bool)

	// Locate all the assets.
	for _, input := range c.Input {
		if err := findFiles(c, input.Path, c.Prefix, input.Recursive, &toc, visitedPaths); err != nil {
			return err
		}
	}

	out, err := os.OpenFile(c.Output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer out.Close()

	return baseTemplate.Execute(out, struct {
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
{{- end -}}

{{- if $.Config.Restore}}

{{template "restore" . -}}
{{- end}}
`))
