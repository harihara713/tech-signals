package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// golang.org/x/net/html
// github.com/PuerkitoBio/goquery
// github.com/gocolly/colly

func main() {
	fmt.Println("Web Scraper")
	args := os.Args
	if len(args) == 1 {
		fmt.Printf("No URL provided!")
		return
	}

	url := args[1]
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching url: %v\n", err)
		return
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Println("Response body:\n", string(data))
}
