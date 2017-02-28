// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

// binAsset holds information about a single asset to be processed.
type binAsset struct {
	Path         string // Relative path.
	AbsPath      string // Absolute path, only used for Debug.
	Name         string // Key used in TOC -- name by which asset is referenced.
	OriginalName string // Original Name before hashing applied to Name.
	Hash         []byte // Generated hash of file.
}

type tocSorter []binAsset

func (s tocSorter) Len() int {
	return len(s)
}

func (s tocSorter) Less(i, j int) bool {
	return s[i].OriginalName < s[j].OriginalName
}

func (s tocSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
