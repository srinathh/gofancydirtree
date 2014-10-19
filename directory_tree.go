// Package gofancydirtree constructs a tree of a supplied root directory in a format that is directly JSON
// serializable into the format expected by the Fancytree JQuery UI Javascript plugin. This makes it simple
// for servers written in Go Language to display advanced directory trees in modern web browsers. See the
// demo folder for an example of the usage for this
//
// Usage
//
// Call the NewTree() function passing the root of the directory tree you want to create. Pass true as the
// second parameter if you want to ignore hidden files and directories that begin with a `.`
//		dirtree, dirmap, err := gofancydirtree.NewTree("/home/me", true)
//
// License
//
// https://github.com/srinathh/gofancydirtree/blob/master/LICENSE
//
// See Also
//
// https://github.com/marcinwyszynski/directory_tree - The Go Package from which this project is forked
//
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
	FullPath string       `json:"-"`                  //The full file path returned by filepath.Walk()
	Info     *os.FileInfo `json:"-"`                  //The os.FileInfo structure returned by filepath.Walk()
	Children []*Node      `json:"children,omitempty"` //The children to this directory tree node
	Title    string       `json:"title"`              //File or Directory name serialized as title
	Key      string       `json:"key"`                //Hex encoded fnv64 hash of the FullPath
	Folder   bool         `json:"folder,omitempty"`   //Is this a folder - ie. Info.IsDir()
}

// Helper function to get a path's parent path (OS-specific).
func getParentPath(path string) string {
	els := strings.Split(path, string(os.PathSeparator))
	return strings.Join(els[:len(els)-1], string(os.PathSeparator))
}

// This creates a directory hierarchy rooted at the specified path and returns a tree structure suitable for
// json serialization for Fancytree plugin. the Nodes contain a key that are also coded in a hashmap returned
// suitable for retrieval in client server communication. The second boolean parameter if true
// skips "hidden" directories and files that typically start with a `.` like `.git`
//
// The first return parameter is the tree structure composed of *Node elements that you can serialize with
// json.Encoder to produce a source to the FancyTree JQuery UI Javascript Plugin. See the demo folder in
// the repository for an example.
//
// The second return parameter dirmap is a map of hash keys (fnv hashed paths as hex strings) to pointers
// to individual nodes in the tree structure. These keys are also set as values for the Node.Key element
// in each node that is emitted as key on JSON serialization. These can be used in Javascript functions
// to communicate events from the client code back tot he server for individual directories or files
// without leaking the absolute directory structure in the server
func NewTree(root string, skiphidden bool) (*Node, map[string]*Node, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, nil, err
	}

	parents := make(map[string]*Node)

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if skiphidden && strings.Index(info.Name(), ".") == 0 {
			return filepath.SkipDir
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
	return hex.EncodeToString(hasher.Sum(nil))
}
