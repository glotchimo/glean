package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"text/template"

	"github.com/fsnotify/fsnotify"
)

var (
	PATH string
	PORT string

	CONF Conf

	//go:embed tmpl/index.html
	INDEX_HTML string
	INDEX_TMPL *template.Template

	//go:embed tmpl/post.html
	POST_HTML string
	POST_TMPL *template.Template
)

func watch(ch chan error) {
	log.Printf("starting FS watcher on %s\n", CONF.PostsPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		ch <- fmt.Errorf("error creating watcher: %w", err)
		return
	}
	defer watcher.Close()

	if err := watcher.Add(CONF.PostsPath); err != nil {
		ch <- fmt.Errorf("error adding posts folder: %w", err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				post, err := MakePost(event.Name)
				if err != nil {
					ch <- fmt.Errorf("error making post: %w", err)
					return
				}

				subject := fmt.Sprintf("P.T: %s", strings.TrimPrefix(strings.TrimSuffix(event.Name, ".md"), CONF.PostsPath))
				content := bytes.Buffer{}
				if err := POST_TMPL.Execute(&content, post); err != nil {
					ch <- fmt.Errorf("error executing template: %w", err)
					return
				}

				if err := SendEmail(subject, content.String()); err != nil {
					ch <- fmt.Errorf("error sending email: %w", err)
					return
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			ch <- fmt.Errorf("error in watcher: %w", err)
			return
		}
	}
}

func serve(ch chan error) {
	log.Printf("starting HTTP server on port %s\n", PORT)

	http.HandleFunc("/", SendIndex)
	http.HandleFunc("/posts/", SendPost)
	http.HandleFunc("/register", TakeEmail)

	ch <- http.ListenAndServe(":"+PORT, nil)
}

func init() {
	flag.StringVar(&PATH, "conf", "~/.config/glean/conf.yml", "Path to configuration")
	flag.StringVar(&PORT, "port", "8080", "Port to listen on")
	flag.Parse()

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
