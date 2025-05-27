# RSS Reader with Headless Browser

This Go program uses a headless Chrome browser to fetch the final HTML content from a URL, following all redirects.

## Prerequisites

- Go 1.16 or higher
- Chrome or Chromium browser installed

## Installation

1. Clone this repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

Run the program:

```bash
go run main.go
```

The program will:
1. Launch a headless Chrome browser
2. Navigate to the specified URL
3. Wait for the page to fully load
4. Print the final HTML content to stdout

## Dependencies

- [chromedp](https://github.com/chromedp/chromedp): Chrome Debugging Protocol client for Go