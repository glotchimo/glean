package main

import (
	_ "embed"
	"log"
	"os"
	"os/signal"
	"text/template"
)

var (
	HOST = os.Getenv("RAILWAY_STATIC_URL")
	PORT = os.Getenv("GLEAN_PORT")
	PATH = os.Getenv("GLEAN_PATH")
	PASS = os.Getenv("GLEAN_PASS")

	TITLE  = os.Getenv("GLEAN_TITLE")
	AUTHOR = os.Getenv("GLEAN_AUTHOR")
	EMAIL  = os.Getenv("GLEAN_EMAIL")

	//go:embed tmpl/index.html
	INDEX_HTML string
	INDEX_TMPL = template.Must(template.New("index").Parse(INDEX_HTML))

	//go:embed tmpl/post.html
	POST_HTML string
	POST_TMPL = template.Must(template.New("post").Parse(POST_HTML))
)

func serve(ch chan error) {
	log.Printf("starting HTTP server on port %s\n", PORT)

	http.HandleFunc("/", SendIndex)
	http.HandleFunc("/new", TakePost)
	http.HandleFunc("/posts/", SendPost)
	http.HandleFunc("/rss", SendFeed)

	ch <- http.ListenAndServe(":"+PORT, nil)
}

func main() {
	serveErrs := make(chan error)
	go serve(serveErrs)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case err := <-serveErrs:
			log.Fatalf("error in server: %s\n", err.Error())
		case sig := <-signals:
			log.Printf("received %s signal, shutting down...\n", sig.String())
			os.Exit(0)
		}
	}
}
