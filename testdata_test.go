// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"

	chacha20rand "github.com/tmthrgd/go-rand"
	"github.com/zach-klippenstein/goregen"
)

type testStat struct{ f *testFile }

func (s testStat) Name() string {
	_, name := path.Split(s.f.path)
	return name
}

func (s testStat) Size() int64        { return s.f.size }
func (s testStat) Mode() os.FileMode  { return s.f.mode }
func (s testStat) ModTime() time.Time { return s.f.time }
func (s testStat) IsDir() bool        { return false }
func (s testStat) Sys() interface{}   { return nil }

type testFile struct {
	path string
	seed [chacha20rand.SeedSize]byte
	size int64
	time time.Time
	mode os.FileMode
}

func (f *testFile) Name() string { return f.path }
func (f *testFile) Path() string { return f.path }

func (f *testFile) AbsolutePath() string {
	return path.Join("/path/to/test/data", f.path)
}

func (f *testFile) Open() (io.ReadCloser, error) {
	r, err := chacha20rand.New(f.seed[:])
	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(io.LimitReader(r, f.size)), nil
}

func (f *testFile) Stat() (os.FileInfo, error) {
	return testStat{f}, nil
}

var (
	numTestFiles = flag.Uint("testfiles", 25, "the number of random test files to add")
	maxFileSize  = flag.Uint("filesize", 512, "the maximum size of random test files")
)

var testFiles Files

func setupTestFiles() {
	mrand := rand.New(rand.NewSource(int64(0)))

	dg, err := regen.NewGenerator("[a-zA-Z0-9_. -][a-zA-Z0-9/_. -]{15,31}", &regen.GeneratorArgs{
		RngSource: mrand,
	})
	if err != nil {
		panic(err)
	}

	var dir string

	testFiles = make(Files, 0, *numTestFiles)
	for i := uint(0); i < *numTestFiles; i++ {
		rand := rand.New(rand.NewSource(int64(i + 1)))

		// regen seems biased towards short names when using
		//  [a-zA-Z0-9_. -]{1,72}
		regex := fmt.Sprintf("[a-zA-Z0-9_. -]{%d}", 1+rand.Intn(72))
		fg, err := regen.NewGenerator(regex, &regen.GeneratorArgs{
			RngSource: rand,
		})
		if err != nil {
			panic(err)
		}

		f := &testFile{
			path: path.Join(dir, fg.Generate()),
			size: int64(rand.Intn(int(*maxFileSize))),
			time: time.Unix(rand.Int63(), rand.Int63()),
			mode: os.FileMode(rand.Intn(int(os.ModePerm))),
		}
		rand.Read(f.seed[:])

		testFiles = append(testFiles, f)

		switch mrand.Intn(5) {
		case 0:
			dir = dg.Generate()
		case 1:
			dir = path.Join(dir, dg.Generate())
		}
	}
}
