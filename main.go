package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://engineering.fb.com/category/open-source/"
	resp, err := parseURL(url)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing the document: %v\n", err)
		return
	}

	doc.Find(".article-grids .row .article-grid").Each(func(i int, s *goquery.Selection) {
		article := s.Find("article")
		var title, url string

		article.Find("a").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				href, _ := s.Attr("href")
				url = href
				return
			}

			title = strings.TrimSpace(s.Text())
		})

		timeStr := article.Find("time").Text()
		t, _ := time.Parse("2006-01-02 15:04:05", timeStr)

		fmt.Printf("Title: %s\nURL: %s\nCreated: %v\n\n", title, url, t)
	})

}

func parseURL(url string) (*http.Response, error) {
	return http.Get(url)
}
