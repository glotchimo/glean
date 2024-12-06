package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"

	"github.com/joho/godotenv"
)

var (
	HOST = os.Getenv("RAILWAY_STATIC_URL")
	PORT = os.Getenv("PORT")
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
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: error loading .env file: %v\n", err)
	}

	// Set default port if not set
	if PORT == "" {
		PORT = "8080"
	}

	// Initialize storage
	if err := initStorage(); err != nil {
		log.Fatalf("Failed to initialize storage: %v\n", err)
	}

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
