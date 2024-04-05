package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

func CrawlWebpage(rootURL string, maxDepth int) ([]string, error) {
	if maxDepth == 0 {
		return nil, nil
	}
	visited := make(map[string]bool)
	queue := SafeUrlQueue{urls: make(map[string]int)}
	baseURL, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}

	visited[rootURL] = true
	queue.Push(rootURL, 0)

	links := []string{rootURL}

	for queue.Len() > 0 {
		link, depth := queue.Pop()
		if depth >= maxDepth {
			continue
		}

		pageLinks, err := getPageLinks(link, baseURL)
		if err != nil {
			return nil, err
		}

		for _, pgLink := range pageLinks {
			if !visited[pgLink] {
				visited[pgLink] = true
				queue.Push(pgLink, depth+1)
				links = append(links, pgLink)
			}
		}
	}

	return links, nil
}

func getPageLinks(link string, baseURL *url.URL) ([]string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	baseURLPath := baseURL.Path

	tokenizer := html.NewTokenizer(resp.Body)
	links := []string{}

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()
		if tokenType == html.StartTagToken && token.Data == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					parsedURL, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}

					if parsedURL.IsAbs() {
						if parsedURL.Host == baseURL.Host || parsedURL.Host == "" {
							links = appendIfNotExists(links, parsedURL.String())
						}
					} else {
						joinedURL := baseURL.Scheme + "://" + baseURL.Host + baseURLPath + parsedURL.Path
						links = appendIfNotExists(links, joinedURL)
					}
				}
			}
		}
	}

	return links, nil
}

func appendIfNotExists(slice []string, item string) []string {
	for _, ele := range slice {
		if ele == item {
			return slice
		}
	}
	return append(slice, item)
}

func main() {
	const (
		defaultURL      = "https://www.example.com/"
		defaultMaxDepth = 3
	)
	urlFlag := flag.String("url", defaultURL, "the url that you want to crawl")
	maxDepth := flag.Int("depth", defaultMaxDepth, "the maximum number of links deep to traverse")
	flag.Parse()

	links, err := CrawlWebpage(*urlFlag, *maxDepth)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	fmt.Println("Links")
	fmt.Println("-----")
	for i, l := range links {
		fmt.Printf("%03d. %s\n", i+1, l)
	}
	fmt.Println()
}
