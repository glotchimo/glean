package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

// ServeIndex generates an index from the files in `posts/`.
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(conf.Path)
	if err != nil {
		http.Error(w, "error reading posts directory", http.StatusInternalServerError)
		return
	}

	idx := Index{Conf: conf}
	for _, e := range files {
		name := e.Name()
		if name[len(name)-3:] != ".md" {
			continue
		}

		name = strings.TrimSuffix(name, ".md")
		idx.Titles = append(idx.Titles, name)
	}

	if err := indexTmpl.Execute(w, idx); err != nil {
		http.Error(w, "error executing template", http.StatusInternalServerError)
	}
}

// ServePost reads from `posts/*.md` and translates to HTML.
func ServePost(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/posts")
	f, err := os.Open(conf.Path + path + ".md")
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

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	post := Post{
		Conf:    conf,
		Content: string(markdown.ToHTML(md, parser, nil)),
	}

	if err := postTmpl.Execute(w, post); err != nil {
		http.Error(w, "error executing template", http.StatusInternalServerError)
	}
}
