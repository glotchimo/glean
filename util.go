package main

import (
	"bufio"
	"fmt"
	"io"
	"net/mail"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

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
