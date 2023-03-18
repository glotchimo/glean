package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
)

// ServeIndex generates an index from the files in `posts/`.
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("posts")
	if err != nil {
		http.Error(w, "error reading posts directory", http.StatusInternalServerError)
		return
	}

	idx := Index{Conf: conf}
	for _, e := range files {
		idx.Titles = append(idx.Titles, e.Name())
	}

	if err := indexTmpl.Execute(w, idx); err != nil {
		http.Error(w, "error executing template", http.StatusInternalServerError)
	}
}

// ServePost reads from `posts/*.md`, translates to HTML.
func ServePost(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("." + r.URL.Path)
	if err != nil {
		http.Error(w, "error opening post file", http.StatusNotFound)
		return
	}
	defer f.Close()

	md, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, "error reading post content", http.StatusInternalServerError)
		return
	}

	html := markdown.ToHTML(md, nil, nil)
	post := Post{Conf: conf, Content: string(html)}

	if err := postTmpl.Execute(w, post); err != nil {
		panic(err)
	}
}
