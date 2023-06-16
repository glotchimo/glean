package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"text/template"
)

var (
	CONF Conf

	//go:embed tmpl/index.html
	INDEX_HTML string
	INDEX_TMPL *template.Template

	//go:embed tmpl/post.html
	POST_HTML string
	POST_TMPL *template.Template
)

func init() {
	if err := LoadConf(); err != nil {
		log.Fatal(err)
	}

	INDEX_TMPL = template.Must(template.New("index").Parse(INDEX_HTML))
	POST_TMPL = template.Must(template.New("post").Parse(POST_HTML))
}

func main() {
	watchErrs := make(chan error)
	serveErrs := make(chan error)

	go watch(watchErrs)
	go serve(serveErrs)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case err := <-watchErrs:
			log.Printf("error in watcher: %s\n", err.Error())
		case err := <-serveErrs:
			log.Fatalf("error in server: %s\n", err.Error())
		case sig := <-signals:
			log.Printf("received %s signal, shutting down...\n", sig.String())
			os.Exit(0)
		}
	}
}
