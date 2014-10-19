// Package directory_tree provides a way to generate a directory tree.
//
// Example usage:
//
//	tree, err := directory_tree.NewTree("/home/me")
//
// I did my best to keep it OS-independent but truth be told I only tested it
// on OS X and Debian Linux so mileage may vary. You've been warned.
package directory_tree

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

// Node represents a node in a directory tree.
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

// Create directory hierarchy.
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
	hash := md5.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}
