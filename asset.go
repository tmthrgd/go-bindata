// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

// binAsset holds information about a single asset to be processed.
type binAsset struct {
	Path string // Full file path.
	Name string // Key used in TOC -- name by which asset is referenced.

	OriginalName string // Original Name before hashing applied to Name.
	Hash         []byte // Generated hash of file.
}
