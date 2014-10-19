package directory_tree

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"text/template"
)

func TestDirectoryTree(t *testing.T) {
	tree, _, err := NewTree("test")
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.Encode(tree)

	demotemplate, err := template.ParseFiles("html/templates/demo.html")
	if err != nil {
		t.Fatal(err)
	}
	//testdata := `[{"title": "Node 1", "key": "1"}, {"title": "Folder 2", "key": "2", "folder": true, "children": [ {"title": "Node 2.1", "key": "3"}, {"title": "Node 2.2", "key": "4"} ]} ]`
	http.Handle("/js/", http.FileServer(http.Dir("html/static")))
	http.Handle("/css/", http.FileServer(http.Dir("html/static")))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		demotemplate.Execute(w, buf.String())
	})
	t.Fatal(http.ListenAndServe("127.0.0.1:8989", nil))
}
