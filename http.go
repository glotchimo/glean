package main

import (
	"log"
	"net/http"
	"strings"
)

type Index struct {
	Meta   Meta
	Titles []string
}

type Post struct {
	Meta    Meta
	Content string
}

func SendFeed(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s (%dB)", r.RemoteAddr, r.Method, r.Host, r.ContentLength)

	if r.Method != "GET" {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	index, err := MakeIndex()
	if err != nil {
		log.Println("error making index:", err.Error())
		http.Error(w, "error making index", http.StatusInternalServerError)
		return
	}

	feed, err := MakeFeed(index)
	if err != nil {
		log.Println("error making feed:", err.Error())
		http.Error(w, "error making feed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	w.Write([]byte(feed))
}

func SendIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s (%dB)", r.RemoteAddr, r.Method, r.Host, r.ContentLength)

	if r.Method != "GET" {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	index, err := MakeIndex()
	if err != nil {
		log.Println("error making index:", err.Error())
		http.Error(w, "error making index", http.StatusInternalServerError)
		return
	}

	if err := INDEX_TMPL.Execute(w, index); err != nil {
		log.Println(err.Error())
		http.Error(w, "error executing template", http.StatusInternalServerError)
	}
}

func SendPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s (%dB)", r.RemoteAddr, r.Method, r.Host, r.ContentLength)

	if r.Method != "GET" {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/posts")
	post, err := MakePost(CONF.PostsPath + path + ".md")
	if err != nil {
		log.Println("error making post:", err.Error())
		http.Error(w, "error making post", http.StatusInternalServerError)
		return
	}

	if err := POST_TMPL.Execute(w, post); err != nil {
		log.Println(err.Error())
		http.Error(w, "error executing template", http.StatusInternalServerError)
	}
}

func TakeEmail(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s (%dB)", r.RemoteAddr, r.Method, r.Host, r.ContentLength)

	if r.Method != "POST" {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	if err := SaveEmail(email); err != nil {
		log.Println("error saving email:", err.Error())
		http.Error(w, "error saving email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
