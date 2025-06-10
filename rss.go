package main

import (
	"github.com/uller91/goGator/internal/database"
	"context"
	"net/http"
	"fmt"
	"encoding/xml"
	"html"
	"io"
	"github.com/lib/pq"
	"time"
	"github.com/google/uuid"
	"database/sql"
	"strings"
	"log"
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

func scrapeFeeds(s *state, ctx context.Context) error {
	nextFeed, err := s.database.GetNextFeedToFetch(ctx)
	if err!= nil {
		if pqError, ok := err.(*pq.Error); ok {
			return pqError
		} else {
			return err
		}
	}

	fmt.Println(nextFeed.Name)
	fmt.Println(nextFeed.LastFetchedAt)

	param := database.MarkFeedFetchedParams{UpdatedAt: time.Now(), ID: nextFeed.ID}
	markedFeed, err := s.database.MarkFeedFetched(ctx, param)
	if err!= nil {
		if pqError, ok := err.(*pq.Error); ok {
			return pqError
		} else {
			return err
		}
	}

	rss, err := fetchFeed(ctx, markedFeed.Url)
	if err!= nil {
			return err
	} 

	//fmt.Printf("Feed %v\n", rss.Channel.Title)
	//fmt.Println("Content:")



	for _, item := range rss.Channel.Item {

		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		param := database.CreatePostParams{
			ID: uuid.New(), 
			CreatedAt: time.Now(), 
			UpdatedAt: time.Now(), 
			Title: item.Title, 
			Url: item.Link, 
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID: markedFeed.ID,
		}

		_, err := s.database.CreatePost(context.Background(), param)
		if err!= nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			} else {
				log.Printf("Couldn't create post: %v", err)
				continue
			}
	}
	}

	log.Printf("Feed %v collected, %v posts found", markedFeed.Name, len(rss.Channel.Item))
	return nil
}
