// Package gofancydirtree constructs a tree of a supplied root directory in a format that is directly JSON
// serializable into the format expected by the Fancytree JQuery UI Javascript plugin. This makes it simple
// for servers written in Go Language to display advanced directory trees in modern web browsers. See the
// demo folder for an example of the usage for this
//
// Usage:
// dirtree, dirmap, err := gofancydirtree.NewTree("/home/me")
//
// Home Page: https://github.com/srinathh/gofancydirtree
// License: https://github.com/srinathh/gofancydirtree/blob/master/LICENSE
//
// See Also:
// https://github.com/marcinwyszynski/directory_tree - The Go Package from which this project is forked
// https://github.com/mar10/fancytree - The JQuery UI Plugin for generating Trees in javascript
package gofancydirtree

import (
	"encoding/hex"
	"hash/fnv"
	"os"
	"path/filepath"
	"strings"
)

// Node represents a node in a directory tree. FullPath and Info are available but not encoded into JSON
type Node struct {
	FullPath string       `json:"-"`
	Info     *os.FileInfo `json:"-"`
	Children []*Node      `json:"children,omitempty"`
	Name     string       `json:"title"`
	Key      string       `json:"key"`
	Folder   bool         `json:"folder,omitempty"`
}

// Helper function to get a path's parent path (OS-specific).
func getParentPath(path string) string {
	els := strings.Split(path, string(os.PathSeparator))
	return strings.Join(els[:len(els)-1], string(os.PathSeparator))
}

// Create directory hierarchy rooted at the specified path and return a tree structure suitable for json serialization
// for Fancytree or for tree traversal and hashmap of keys to each node in the tree suitable retreival during client
// server communication.

// The first return parameter (dirtree) is the tree structure composed of *Node elements that you can
// straightaway serialize with json.encoder to produce a dataset that can be fed as source to the FancyTree
// JQuery UI Plugin. See the demo folder in the repository for an example.
//
// The second return parameter (dirmap) is a map of hash keys to individual nodes in the dirtree. these are
// set as values for the "key" element in the Fancytree json dataset and can be used for event communication
// between the client side javascript and the server without leaking the absolute directory structure in
// the server
func NewTree(root string) (*Node, map[string]*Node, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, nil, err
	}

	parents := make(map[string]*Node)

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		key := doHash(path)
		parents[key] = &Node{path, &info, []*Node{}, info.Name(), key, info.IsDir()}
		return nil
	}

	if err = filepath.Walk(absRoot, walkFunc); err != nil {
		return nil, nil, err
	}

	var result *Node = nil

	for _, node := range parents {
		parentPath := getParentPath(node.FullPath)
		parent, exists := parents[doHash(parentPath)]
		if !exists { // If a parent does not exist, this is the root.
			result = node
		} else {
			parent.Children = append(parent.Children, node)
		}
	}
	return result, parents, nil
}

func doHash(s string) string {
	hasher := fnv.New64()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}
