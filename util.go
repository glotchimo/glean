package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
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

				subject := fmt.Sprintf(
					"%s: %s",
					CONF.Meta.Title,
					strings.TrimPrefix(strings.TrimSuffix(event.Name, ".md"), CONF.PostsPath+"/"))
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
	http.HandleFunc("/rss", SendFeed)
	http.HandleFunc("/posts/", SendPost)
	http.HandleFunc("/editor", SendEditor)
	http.HandleFunc("/post", TakePost)
	http.HandleFunc("/register", TakeEmail)

	ch <- http.ListenAndServe(":"+PORT, nil)
}

func MakeIndex() (*Index, error) {
	files, err := os.ReadDir(CONF.PostsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading posts directory: %w", err)
	}

	index := Index{Meta: CONF.Meta}
	for _, e := range files {
		name := e.Name()
		if name[len(name)-3:] != ".md" {
			continue
		}

		name = strings.TrimSuffix(name, ".md")
		index.Titles = append(index.Titles, name)
	}

	for i := len(index.Titles)/2 - 1; i >= 0; i-- {
		j := len(index.Titles) - 1 - i
		index.Titles[i], index.Titles[j] = index.Titles[j], index.Titles[i]
	}

	return &index, nil
}

func MakePost(path string) (*Post, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening post file: %w", err)
	}
	defer f.Close()

	md, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading post content: %w", err)
	}

	ext := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(ext)
	post := Post{
		Meta:    CONF.Meta,
		Content: string(markdown.ToHTML(md, parser, nil)),
	}

	return &post, nil
}

func SaveEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("error parsing email: %w", err)
	}

	emails, err := ReadEmails()
	if err != nil {
		return fmt.Errorf("error reading emails: %w", err)
	}

	unique := map[string]bool{email: true}
	for _, e := range emails {
		unique[e] = true
	}

	f, err := os.OpenFile(CONF.EmailsPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening emails file: %w", err)
	}
	defer f.Close()

	for e := range unique {
		if _, err := f.WriteString(e + "\n"); err != nil {
			return fmt.Errorf("error saving email: %w", err)
		}
	}

	return nil
}

func ReadEmails() ([]string, error) {
	f, err := os.Open(CONF.EmailsPath)
	if err != nil {
		return nil, fmt.Errorf("error opening emails file: %w", err)
	}
	defer f.Close()

	unique := map[string]bool{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		unique[s.Text()] = true
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("error scanning emails file: %w", err)
	}

	emails := []string{}
	for e := range unique {
		emails = append(emails, e)
	}

	return emails, nil
}
