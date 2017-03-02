// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package chain

import (
	"os"
	"path/filepath"
)

// AssetAndInfo represents the generated AssetAndInfo method.
type AssetAndInfo func(name string) (data []byte, info os.FileInfo, err error)

// AssetAndInfoChain represents a chain of AssetAndInfo methods
// that will be called in turn.
type AssetAndInfoChain []AssetAndInfo

// AssetAndInfo loads and returns the asset and asset info for the
// given name. It returns an error if the asset could not be found
// or could not be loaded.
//
// It tries each AssetAndInfo in the chain and returns the first
// successful result or the first error that is not os.ErrNotExist.
func (ch AssetAndInfoChain) AssetAndInfo(name string) (data []byte, info os.FileInfo, err error) {
	path := filepath.ToSlash(name)

	for _, fn := range ch {
		data, info, err = fn(path)
		if err == nil || !os.IsNotExist(err) {
			return
		}
	}

	return nil, nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}
