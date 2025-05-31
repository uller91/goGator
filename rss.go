package main

import (
	"context"
	"net/http"
	"fmt"
	"encoding/xml"
	"html"
	"io"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
        return nil, fmt.Errorf("error making a new request: %w", err)
    }

	req.Header.Set("User-Agent","gator")
	
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
        return nil, fmt.Errorf("error receiving response: %w", err)
    }
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response: %w", err)
    }
	
	xmlItems := []RSSItem{}
	var xmlData RSSFeed
	xmlData.Channel.Item = xmlItems
	//xmlData := RSSFeed{Channel struct{Title: "", Link: "", Description: "", Item: xmlItems,},}
	if err := xml.Unmarshal(data, &xmlData); err != nil {
        return nil, fmt.Errorf("error unmarshaling xml: %w", err)
    }
	
	xmlData.Channel.Title = html.UnescapeString(xmlData.Channel.Title)
	xmlData.Channel.Description = html.UnescapeString(xmlData.Channel.Description)
	for _, item := range xmlData.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &xmlData, nil
}