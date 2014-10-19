package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/srinathh/gofancydirtree"
	"net/http"
	"text/template"
)

func main() {
	tree, _, err := gofancydirtree.NewTree("test", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.Encode(tree)

	demotemplate, err := template.ParseFiles("html/templates/demo.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Handle("/js/", http.FileServer(http.Dir("html/static")))
	http.Handle("/css/", http.FileServer(http.Dir("html/static")))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		demotemplate.Execute(w, buf.String())
	})
	fmt.Println("Serving test server on 127.0.0.1:8989")
	fmt.Println(http.ListenAndServe("127.0.0.1:8989", nil))
}
