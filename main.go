package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/markusmobius/go-trafilatura"
)

// OutputData はJSON出力用の構造体
type OutputData struct {
	Title         string    `json:"title,omitempty"`
	Author        string    `json:"author,omitempty"`
	URL           string    `json:"url,omitempty"`
	Hostname      string    `json:"hostname,omitempty"`
	Description   string    `json:"description,omitempty"`
	Sitename      string    `json:"sitename,omitempty"`
	Date          time.Time `json:"date,omitempty"` // time.Time 型に変更
	Categories    []string  `json:"categories,omitempty"`
	Tags          []string  `json:"tags,omitempty"`
	ID            string    `json:"id,omitempty"`
	Fingerprint   string    `json:"fingerprint,omitempty"`
	License       string    `json:"license,omitempty"`
	Language      string    `json:"language,omitempty"`
	Image         string    `json:"image,omitempty"` // Metadata.Image を使用
	PageType      string    `json:"page_type,omitempty"`
	ContentText   string    `json:"content_text,omitempty"`
	CommentsText  string    `json:"comments_text,omitempty"`
	// Note: Multiple images and links are not directly available as simple lists in ExtractResult.
	// Further processing of ContentNode or specific options might be needed.
}

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
		// JSON出力用のデータ構造に詰め替える
		output := OutputData{
			ContentText:  extracted.ContentText,
			CommentsText: extracted.CommentsText,

			// extracted.Metadata はポインタではないため、nilチェックは不要
			Title:       extracted.Metadata.Title,
			Author:      extracted.Metadata.Author,
			URL:         extracted.Metadata.URL,
			Hostname:    extracted.Metadata.Hostname,
			Description: extracted.Metadata.Description,
			Sitename:    extracted.Metadata.Sitename,
			Date:        extracted.Metadata.Date, // 直接代入
			Categories:  extracted.Metadata.Categories,
			Tags:        extracted.Metadata.Tags,
			ID:          extracted.Metadata.ID,
			Fingerprint: extracted.Metadata.Fingerprint,
			License:     extracted.Metadata.License,
			Language:    extracted.Metadata.Language,
			Image:       extracted.Metadata.Image,
			PageType:    extracted.Metadata.PageType,
		}

		// 抽出結果をJSONにマーシャリング
		jsonData, err := json.MarshalIndent(output, "", "  ") // インデントして見やすくする
		if err != nil {
			log.Fatalf("Error marshalling to JSON: %v", err)
		}

		// JSONデータをファイルに書き込む
		fileName := "extracted_content.json"
		err = os.WriteFile(fileName, jsonData, 0644)
		if err != nil {
			log.Fatalf("Error writing JSON to file: %v", err)
		}
		fmt.Printf("Extracted content saved to %s\n", fileName)
	} else {
		fmt.Println("No content extracted.")
	}
}
