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

type Google struct {
}

func (g Google) Name() string {
	return "Google"
}

func (g Google) Fetch(as *article.ArticleStore) error {
	url := "https://developers.googleblog.com/search/?technology_categories=Web"
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Google: Error getting response", "err", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Google: Error response", "status", resp.StatusCode)
		return fmt.Errorf("Google blogs response status: '%d'", resp.StatusCode)
	}
	// extract the articles details
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		slog.Error("Google: Error parsing the page", "err", err)
		return err
	}

	doc.Find(".search-results__results-wrapper ul").Children().Each(func(i int, s *goquery.Selection) {
		art := article.Article{}

		div := s.Find(".search-result__wrapper div")
		pText := div.Find("p").Text()
		split := strings.Split(pText, "/")

		if len(split) >= 2 {
			art.Tags = append(art.Tags, split[1])
		}

		// split[0] is time string before parsing the time make it in the 'Jan 02, 2006' format
		t, err := time.Parse("Jan 02, 2006", split[0])
		if err != nil {
			slog.Error("Google: Error Parsing time", "err", err)
			return
		}

		header := div.Find("h3")
		art.Title = header.Text()
		art.URL = "https://developers.googleblog.com/" + header.Find("a").AttrOr("href", "")

		art.Summary = div.Find(".search-result__summary").Text()
		art.PublishedAt = t
		art.Source = g.Name()
		*as = append(*as, art)
	})

	return nil
}
