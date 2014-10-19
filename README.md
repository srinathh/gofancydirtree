GoFancyDirTree
==============
Package gofancydirtree constructs a tree of a supplied root directory in a format that is directly JSON
serializable into the format expected by the Fancytree JQuery UI Javascript plugin. This makes it simple
for servers written in Go Language to display advanced directory trees in modern web browsers. See the
demo folder for an example of the usage for this

Installation
------------
    go get github.com/srinathh/gofancydirtree

Usage
-----
    dirtree, dirmap, err := gofancydirtree.NewTree("/home/me", true)

This creates a directory hierarchy rooted at the specified path and returns a tree structure suitable 
for json serialization for Fancytree plugin. the Nodes contain a key that are also coded in a hashmap
returned suitable for retrieval in client server communication. The second boolean parameter if true
skips "hidden" directories and files that typically start with a `.` like `.git`

The first return parameter `dirtree` is the tree structure composed of `*Node` elements that you can
serialize with json.Encoder to produce a `source` to the FancyTree JQuery UI Javascript Plugin. See 
the demo folder in the repository for an example.

The second return parameter `dirmap` is a map of hash keys (fnv hashed paths as hex strings) to pointers
to individual nodes in the `dirtree` (ie. `map[string]*Node`). These keys are also set as values for the
`Node.Key` element in each node that is emitted as `key` on JSON serialization. These can be used in
Javascript functions to communicate events from the client code back tot he server for individual directories
or files without leaking the absolute directory structure in the server

See Also
--------
- The [marcinwyszynski/directory_tree](https://github.com/marcinwyszynski/directory_tree) package from which this is forked
- The [Fancytree](https://github.com/mar10/fancytree) JQuery UI Plugin for generating Trees in javascript
