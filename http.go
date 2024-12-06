package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Meta struct {
	Title  string
	Author string
	Email  string
}

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

	path := strings.TrimPrefix(r.URL.Path, "/posts/")
	post, err := MakePost(path)
	if err != nil {
		log.Println("error making post:", err.Error())
		http.Error(w, "error making post", http.StatusInternalServerError)
		return
	}

	if err := POST_TMPL.Execute(w, post); err != nil {
		log.Println("error executing template:", err.Error())
		http.Error(w, "error executing template", http.StatusInternalServerError)
	}
}

func TakePost(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s (%dB)", r.RemoteAddr, r.Method, r.Host, r.ContentLength)

	if r.Method != "POST" {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if !IsKeyValid(r.Header.Get("Authorization")) {
		log.Println("invalid key")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		log.Println("empty title")
		http.Error(w, "empty title", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		log.Println("empty content")
		http.Error(w, "empty content", http.StatusBadRequest)
		return
	}

	name := fmt.Sprintf("%s %s", time.Now().Format(time.DateOnly), title)
	if err := savePost(name, []byte(content)); err != nil {
		log.Println("error saving post:", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
