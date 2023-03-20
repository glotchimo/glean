package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

var (
	conf Conf

	//go:embed tmpl/index.html
	indexHTML string
	indexTmpl *template.Template

	//go:embed tmpl/post.html
	postHTML string
	postTmpl *template.Template
)

type Conf struct {
	Title  string `yaml:"title"`
	Author string `yaml:"author"`
	Email  string `yaml:"email"`
	Link   string `yaml:"link"`
}

type Index struct {
	Conf   Conf
	Titles []string
}

type Post struct {
	Conf    Conf
	Content string
}

func init() {
	f, err := os.Open("conf.yml")
	if err != nil {
		log.Fatal("error opening config:", err.Error())
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&conf); err != nil {
		log.Fatal("error parsing config:", err.Error())
	}

	indexTmpl = template.Must(template.New("index").Parse(indexHTML))
	postTmpl = template.Must(template.New("post").Parse(postHTML))
}

func main() {
	http.HandleFunc("/", ServeIndex)
	http.HandleFunc("/posts/", ServePost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
