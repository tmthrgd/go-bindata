// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type assetTree struct {
	Asset    binAsset
	Children map[string]*assetTree
}

func newAssetTree() *assetTree {
	tree := &assetTree{}
	tree.Children = make(map[string]*assetTree)
	return tree
}

func (node *assetTree) child(name string) *assetTree {
	rv, ok := node.Children[name]
	if !ok {
		rv = newAssetTree()
		node.Children[name] = rv
	}
	return rv
}

func (root *assetTree) Add(route []string, asset binAsset) {
	for _, name := range route {
		root = root.child(name)
	}
	root.Asset = asset
}

func ident(w io.Writer, n int) {
	for i := 0; i < n; i++ {
		w.Write([]byte{'\t'})
	}
}

func (root *assetTree) funcOrNil() string {
	if root.Asset.Func == "" {
		return "nil"
	} else {
		return root.Asset.Func
	}
}

func (root *assetTree) writeGoMap(w io.Writer, nident int) {
	if nident == 0 {
		io.WriteString(w, "&bintree")
	}
	fmt.Fprintf(w, "{%s, map[string]*bintree{", root.funcOrNil())

	if len(root.Children) > 0 {
		io.WriteString(w, "\n")

		// Sort to make output stable between invocations
		filenames := make([]string, len(root.Children))
		i := 0
		for filename, _ := range root.Children {
			filenames[i] = filename
			i++
		}
		sort.Strings(filenames)

		for _, p := range filenames {
			ident(w, nident+1)
			fmt.Fprintf(w, `"%s": `, p)
			root.Children[p].writeGoMap(w, nident+1)
		}
		ident(w, nident)
	}

	io.WriteString(w, "}}")
	if nident > 0 {
		io.WriteString(w, ",")
	}
	io.WriteString(w, "\n")
}

func (root *assetTree) WriteAsGoMap(w io.Writer) error {
	_, err := fmt.Fprint(w, `type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = `)
	root.writeGoMap(w, 0)
	return err
}

func writeTree(w io.Writer, toc []binAsset) error {
	_, err := io.WriteString(w, `// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
			}
		}
	}
	if node.Func != nil {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

`)
	if err != nil {
		return err
	}
	tree := newAssetTree()
	for i := range toc {
		pathList := strings.Split(toc[i].Name, "/")
		tree.Add(pathList, toc[i])
	}
	return tree.WriteAsGoMap(w)
}
