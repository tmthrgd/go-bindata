// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"io"
	"strings"
	"text/template"
)

// Generator generates Go code that embeds static assets.
type Generator struct {
	c       Config
	toc     []binAsset
	visited map[string]struct{}
}

// New returns a new Generator with a given configuration.
func New(c *Config) (*Generator, error) {
	g := &Generator{
		c:       *c,
		visited: make(map[string]struct{}),
	}

	// Ensure our configuration has sane values.
	if err := g.c.validate(); err != nil {
		return nil, err
	}

	return g, nil
}

// FindFiles adds all files inside a directory to the
// generated output. If recursive is true, files within
// subdirectories of path will also be included.
func (g *Generator) FindFiles(path string, recursive bool) error {
	return findFiles(&g.c, path, g.c.Prefix, recursive, &g.toc, g.visited)
}

// WriteTo writes the generated Go code to w.
func (g *Generator) WriteTo(w io.Writer) (n int64, err error) {
	lw := lenWriter{W: w}
	err = baseTemplate.Execute(&lw, struct {
		Config    *Config
		AssetName bool
		Assets    []binAsset
	}{&g.c, g.c.HashFormat != NoHash && g.c.HashFormat != NameUnchanged, g.toc})
	return lw.N, err
}

type lenWriter struct {
	W io.Writer
	N int64
}

func (lw *lenWriter) Write(p []byte) (n int, err error) {
	n, err = lw.W.Write(p)
	lw.N += int64(n)
	return
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
