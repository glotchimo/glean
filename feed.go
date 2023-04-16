package main

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	LastBuildDate string `xml:"lastBuildDate"`
	Items         []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func MakeFeed(index *Index) (string, error) {
	feed := Feed{
		Version: "2.0",
		Channel: Channel{
			Title:         CONF.Meta.Title,
			Link:          CONF.Host,
			LastBuildDate: time.Now().Format(time.RFC1123),
		},
	}

	for _, t := range index.Titles {
		pub, err := time.Parse("2006-01-02", strings.Split(t, " ")[0])
		if err != nil {
			return "", fmt.Errorf("error parsing post timestamp: %w", err)
		}

		item := Item{
			Title:   t,
			Link:    CONF.Host + "/posts/" + t,
			PubDate: pub.Format(time.RFC1123),
		}

		feed.Channel.Items = append(feed.Channel.Items, item)
	}

	out, err := xml.MarshalIndent(feed, "", "\t")
	if err != nil {
		return "", fmt.Errorf("error marshalling feed: %w", err)
	}

	return xml.Header + string(out), nil
}
