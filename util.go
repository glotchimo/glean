package main

import (
	"crypto/subtle"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func IsKeyValid(key string) bool {
	return subtle.ConstantTimeCompare([]byte(key), []byte(PASS)) == 1
}

func MakeIndex() (*Index, error) {
	files, err := os.ReadDir(PATH)
	if err != nil {
		return nil, fmt.Errorf("error reading posts directory: %w", err)
	}

	index := Index{Meta: Meta{TITLE, AUTHOR, EMAIL}}
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
		Meta:    Meta{TITLE, AUTHOR, EMAIL},
		Content: string(markdown.ToHTML(md, parser, nil)),
	}

	return &post, nil
}
