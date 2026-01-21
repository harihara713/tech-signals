package source

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/harihara713/tech-signals/article"
)

type Meta struct {
}

func (m Meta) Name() string {
	return "Meta"
}

func (m Meta) Fetch(as *article.ArticleStore) error {
	url := "https://engineering.fb.com/category/open-source/"
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Meta: Error getting response", "err", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Meta: Error response", "status", resp.StatusCode)
		return fmt.Errorf("Meta blogs response status: '%d'", resp.StatusCode)
	}
	// extract the articles details
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		slog.Error("Meta: Error parsing the page", "err", err)
		return err
	}

	doc.Find(".article-grids .row .article-grid").Each(func(i int, s *goquery.Selection) {
		art := article.Article{}

		articleTag := s.Find("article")
		articleTag.Find("a").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				href, _ := s.Attr("href")
				art.URL = href
				return
			}

			art.Title = strings.TrimSpace(s.Text())
		})

		timeStr := articleTag.Find("time").Text()
		//TODO: Transform the data in Jan 02, 2006 format first
		t, err := time.Parse("Jan 02, 2006", timeStr)
		if err != nil {
			slog.Error("Meta: Error Parsing time", "err", err)
			return
		}

		art.PublishedAt = t
		art.Source = m.Name()

		*as = append(*as, art)
	})

	return nil
}
