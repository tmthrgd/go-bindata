// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"go/parser"
	"go/token"
)

func formatTemplate(name string, data interface{}) (string, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString("package main;")

	if err := baseTemplate.ExecuteTemplate(buf, name, data); err != nil {
		return "", err
	}

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", buf, parser.ParseComments)
	if err != nil {
		return "", err
	}

	buf.Reset()

	if err = printerConfig.Fprint(buf, fset, f); err != nil {
		return "", err
	}

	out := string(bytes.TrimSpace(buf.Bytes()[len("package main\n"):]))
	buf.Reset()
	bufPool.Put(buf)
	return out, nil
}
