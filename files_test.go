// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"strings"
	"testing"
)

func TestFindFiles(t *testing.T) {
	var toc []binAsset
	var knownFuncs = make(map[string]int)
	var visitedPaths = make(map[string]bool)
	err := findFiles(new(Config), "testdata/dupname", "testdata/dupname", true, &toc, knownFuncs, visitedPaths)
	if err != nil {
		t.Errorf("expected to be no error: %+v", err)
	}
	if toc[0].Func == toc[1].Func {
		t.Errorf("name collision")
	}
}

func TestFindFilesWithSymlinks(t *testing.T) {
	var tocSrc []binAsset
	var tocTarget []binAsset

	var knownFuncs = make(map[string]int)
	var visitedPaths = make(map[string]bool)
	err := findFiles(new(Config), "testdata/symlinkSrc", "testdata/symlinkSrc", true, &tocSrc, knownFuncs, visitedPaths)
	if err != nil {
		t.Errorf("expected to be no error: %+v", err)
	}

	knownFuncs = make(map[string]int)
	visitedPaths = make(map[string]bool)
	err = findFiles(new(Config), "testdata/symlinkParent", "testdata/symlinkParent", true, &tocTarget, knownFuncs, visitedPaths)
	if err != nil {
		t.Errorf("expected to be no error: %+v", err)
	}

	if len(tocSrc) != len(tocTarget) {
		t.Errorf("Symlink source and target should have the same number of assets.  Expected %d got %d", len(tocTarget), len(tocSrc))
	} else {
		for i, _ := range tocSrc {
			targetFunc := strings.TrimPrefix(tocTarget[i].Func, "symlinktarget")
			targetFunc = strings.ToLower(targetFunc[:1]) + targetFunc[1:]
			if tocSrc[i].Func != targetFunc {
				t.Errorf("Symlink source and target produced different function lists.  Expected %s to be %s", targetFunc, tocSrc[i].Func)
			}
		}
	}
}

func TestFindFilesWithRecursiveSymlinks(t *testing.T) {
	var toc []binAsset

	var knownFuncs = make(map[string]int)
	var visitedPaths = make(map[string]bool)
	err := findFiles(new(Config), "testdata/symlinkRecursiveParent", "testdata/symlinkRecursiveParent", true, &toc, knownFuncs, visitedPaths)
	if err != nil {
		t.Errorf("expected to be no error: %+v", err)
	}

	if len(toc) != 1 {
		t.Errorf("Only one asset should have been found.  Got %d: %v", len(toc), toc)
	}
}

func TestFindFilesWithSymlinkedFile(t *testing.T) {
	var toc []binAsset

	var knownFuncs = make(map[string]int)
	var visitedPaths = make(map[string]bool)
	err := findFiles(new(Config), "testdata/symlinkFile", "testdata/symlinkFile", true, &toc, knownFuncs, visitedPaths)
	if err != nil {
		t.Errorf("expected to be no error: %+v", err)
	}

	if len(toc) != 1 {
		t.Errorf("Only one asset should have been found.  Got %d: %v", len(toc), toc)
	}
}
