package main

import (
	"crypto/subtle"
	"fmt"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func IsKeyValid(key string) bool {
	return subtle.ConstantTimeCompare([]byte(key), []byte(PASS)) == 1
}

func MakeIndex() (*Index, error) {
	titles, err := listPosts()
	if err != nil {
		return nil, fmt.Errorf("error listing posts: %w", err)
	}

	index := Index{
		Meta:   Meta{TITLE, AUTHOR, EMAIL},
		Titles: titles,
	}

	return &index, nil
}

func MakePost(name string) (*Post, error) {
	content, err := getPost(name)
	if err != nil {
		return nil, fmt.Errorf("error getting post: %w", err)
	}

	ext := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(ext)
	post := Post{
		Meta:    Meta{TITLE, AUTHOR, EMAIL},
		Content: string(markdown.ToHTML(content, parser, nil)),
	}

	return &post, nil
}
