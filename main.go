package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/markusmobius/go-trafilatura"
)

func main() {
	urlString := "https://news.google.com/rss/articles/CBMidkFVX3lxTFBGX1BCaFdpRmVGYkpmX3R4dGU3MG1rVEJVdXhMQ0NfZ1J3LWlNa2JCT0Q0SWJSZHVneVA0NlN6TkJsTnVITVJXSGh2elhBTGJXWENySlphX3RfUExRSVplbUtjVUNqRkpTcU1qVWVvemI1T1BVZkE?oc=5"

	parsedURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	// Brave Browserの実行可能ファイルのパス (macOSの一般的なパス)
	bravePath := "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser"

	// chromedpのExecAllocatorオプションを設定
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // ヘッドレスモードを無効にする
		// chromedp.Flag("disable-gpu", true), // GPUを無効にする場合 (問題が発生する場合)
		chromedp.ExecPath(bravePath),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	defer cancelAlloc()

	// chromedpのコンテキストを作成
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// 必要に応じてタイムアウトを設定 (例: 30秒)
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var htmlContent string
	log.Println("Fetching URL with chromedp...")
	err = chromedp.Run(ctx,
		chromedp.Navigate(urlString),
		// ページ遷移やJSによるコンテンツ読み込みのために5秒待機
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("Waiting for 5 seconds after navigation...")
			time.Sleep(5 * time.Second)
			return nil
		}),
		// 必要であれば、特定の要素が表示されるまで待機するなどのアクションを追加
		// chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.OuterHTML("html", &htmlContent), // ページ全体のHTMLを取得
	)
	if err != nil {
		log.Fatalf("Chromedp error: %v", err)
	}
	log.Println("HTML fetched successfully with chromedp.")

	// 取得したHTMLをio.Readerに変換
	htmlReader := strings.NewReader(htmlContent)

	// Extract content from the fetched HTML with options
	options := trafilatura.Options{
		Focus:          trafilatura.FavorRecall,
		EnableFallback: true,
		OriginalURL:    parsedURL,
	}
	extracted, err := trafilatura.Extract(htmlReader, options)
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
