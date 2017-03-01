// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

const stubFsRoot = "/path/to/test/data"

var testPaths = map[string]*FindFilesOptions{
	"testdata":               {Recursive: true},
	"testdata/ab6.bin":       {Prefix: "testdata"},
	"testdata/ogqS":          {Prefix: "testdata"},
	"testdata/ogqS/qsDM.bin": {Prefix: "testdata/ogqS"},
}

var modTime = time.Unix(123456789, 987654321)

func testStubFileSystem() error {
	fs := afero.NewMemMapFs()

	var fe firstError
	fe.Set(filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return fs.Mkdir(path, info.Mode())
		}

		var fe firstError

		data, err := ioutil.ReadFile(path)
		fe.Set(err)
		fe.Set(afero.WriteFile(fs, path, data, info.Mode()))
		fe.Set(fs.Chtimes(path, modTime, modTime))
		return fe.Err
	}))

	af := afero.Afero{Fs: fs}

	// for testing: path/filepath
	abs = func(path string) (string, error) {
		if filepath.IsAbs(path) {
			return path, nil
		}

		return filepath.Join(stubFsRoot, path), nil
	}
	walk = af.Walk

	// for testing: os
	open = func(name string) (file, error) {
		// This is ok for our use, but beware of:
		// http://spf13.com/post/when-nil-is-not-nil/
		return af.Open(name)
	}
	stat = af.Stat

	return fe.Err
}

type firstError struct {
	Err error
}

func (e *firstError) Set(err error) {
	if e.Err == nil {
		e.Err = err
	}
}
