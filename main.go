package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// URL to fetch
	url := "https://news.google.com/rss/articles/CBMidkFVX3lxTFBGX1BCaFdpRmVGYkpmX3R4dGU3MG1rVEJVdXhMQ0NfZ1J3LWlNa2JCT0Q0SWJSZHVneVA0NlN6TkJsTnVITVJXSGh2elhBTGJXWENySlphX3RfUExRSVplbUtjVUNqRkpTcU1qVWVvemI1T1BVZkE?oc=5"

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create a new Chrome instance
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.ExecPath("/Applications/Brave Browser.app/Contents/MacOS/Brave Browser"),
		chromedp.UserAgent(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36`),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Create a new browser instance
	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Navigate to the URL and wait for the page to load
	var htmlContent string
	var currentURL string
	err := chromedp.Run(browserCtx,
		chromedp.Navigate(url),
		// Wait for the page to be fully loaded
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Location(&currentURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Printf("Navigated to URL: %s", currentURL)
			return nil
		}),
		chromedp.Sleep(5 * time.Second),
		// Get the HTML content
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatalf("Failed to navigate to URL: %v", err)
	}

	// Save the HTML content to a file
	err = os.WriteFile("output.html", []byte(htmlContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write HTML to file: %v", err)
	}

	fmt.Println("HTML content saved to output.html")
}
