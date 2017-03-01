// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing/quick"
	"time"

	"github.com/spf13/afero"
	"github.com/zach-klippenstein/goregen"
)

var testPaths = map[string]*FindFilesOptions{
	"/": {Recursive: true},
}

func testRandomName(rand *rand.Rand) string {
	g, err := regen.NewGenerator(fmt.Sprintf("[a-zA-Z0-9_.-]{%d}", 1+rand.Intn(64-1)), &regen.GeneratorArgs{
		RngSource: rand,
	})
	if err != nil {
		panic(err)
	}

	return g.Generate()
}

type testFileName string

func (testFileName) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(testFileName(testRandomName(rand)))
}

type testFileData []byte

func (testFileData) Generate(rand *rand.Rand, size int) reflect.Value {
	v := make([]byte, 1+rand.Intn(512-1))
	rand.Read(v)
	return reflect.ValueOf(testFileData(v))
}

type testFileModTime time.Time

func (testFileModTime) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(testFileModTime(time.Unix(rand.Int63(), rand.Int63())))
}

type testFileMode os.FileMode

func (testFileMode) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(testFileMode(rand.Intn(int(os.ModePerm))))
}

type testFileMap map[testFileName]struct {
	Data    testFileData
	ModTime testFileModTime
	Mode    testFileMode
}

func testPopulateDirectory(fs afero.Fs, base string, rand *rand.Rand) error {
	v, ok := quick.Value(reflect.TypeOf(testFileMap{}), rand)
	if !ok {
		panic("quick.Value failed")
	}

	var fe firstError

	for name, file := range v.Interface().(testFileMap) {
		path := filepath.FromSlash(filepath.Join(base, string(name)))

		fe.Set(afero.WriteFile(fs, path, file.Data, os.FileMode(file.Mode)))
		fe.Set(fs.Chtimes(path, time.Time(file.ModTime), time.Time(file.ModTime)))
	}

	return fe.Err
}

var fs = afero.NewMemMapFs()

func testStubFileSystem() error {
	rand := rand.New(rand.NewSource(0))

	var fe firstError

	dirName := path.Join("/", testRandomName(rand))
	testPaths[dirName] = &FindFilesOptions{
		Prefix: dirName,
	}

	for _, dir := range [...]string{
		dirName,
		path.Join("/", testRandomName(rand)),
	} {
		fe.Set(testPopulateDirectory(fs, dir, rand))
	}

	af := afero.Afero{Fs: fs}

	// for testing: path/filepath
	abs = func(path string) (string, error) {
		if filepath.IsAbs(path) {
			return path, nil
		}

		return filepath.Join("/", path), nil
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
