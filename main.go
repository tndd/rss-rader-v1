package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/markusmobius/go-trafilatura"
)

func main() {
	urlString := "https://news.google.com/rss/articles/CBMidkFVX3lxTFBGX1BCaFdpRmVGYkpmX3R4dGU3MG1rVEJVdXhMQ0NfZ1J3LWlNa2JCT0Q0SWJSZHVneVA0NlN6TkJsTnVITVJXSGh2elhBTGJXWENySlphX3RfUExRSVplbUtjVUNqRkpTcU1qVWVvemI1T1BVZkE?oc=5"

	// URLをパース
	parsedURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	// Create an HTTP client with a timeout
	httpClient := &http.Client{Timeout: 30 * time.Second}

	// Fetch the URL content
	resp, err := httpClient.Get(urlString)
	if err != nil {
		log.Fatalf("Error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	// Extract content from the response body with options
	options := trafilatura.Options{
		Focus:          trafilatura.FavorRecall,
		EnableFallback: true,
		OriginalURL:    parsedURL,
	}
	extracted, err := trafilatura.Extract(resp.Body, options)
	if err != nil {
		log.Fatalf("Error extracting content: %v", err)
	}

	if extracted != nil {
		fmt.Println("Extracted Text:")
		fmt.Println(extracted.ContentText)
	} else {
		fmt.Println("No content extracted.")
	}
}
