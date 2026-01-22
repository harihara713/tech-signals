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

type Amazon struct {
}

func (a Amazon) Name() string {
	return "Amazon"
}

func (a Amazon) Fetch(as *article.ArticleStore) error {
	url := "https://aws.amazon.com/blogs/architecture/"
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Amazon: Error getting response", "err", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Amazon: Error response", "status", resp.StatusCode)
		return fmt.Errorf("Amazon blogs response status: '%d'", resp.StatusCode)
	}
	// extract the articles details
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		slog.Error("Amazon: Error parsing the page", "err", err)
		return err
	}

	doc.Find("main article").Each(func(i int, s *goquery.Selection) {
		art := article.Article{}

		div := s.Find("div")
		firstChild := div.Find("div").First()
		secondChild := div.ChildrenFiltered("div").Eq(1)

		art.URL = firstChild.Find("a").AttrOr("href", "")
		art.Title = secondChild.Find("h2").Text()
		art.Summary = secondChild.Find("section p").Text()

		var authors []string
		secondChild.Find("footer span span").Each(func(i int, s *goquery.Selection) {
			authors = append(authors, s.Text())
		})

		art.Author = strings.Join(authors, ", ")
		timeVal := secondChild.Find("footer span time").Text()
		timeParts := strings.Fields(timeVal)

		timeStr, err := transformToTimeString(fmt.Sprintf("%s %s, %s", timeParts[1], timeParts[0], timeParts[2]))
		if err != nil {
			slog.Error("Amazon: Error transform to time string", "err", err)
			return
		}

		t, err := time.Parse("Jan 02, 2006", timeStr)
		if err != nil {
			slog.Error("Amazon: Error Parsing time", "err", err)
			return
		}

		art.PublishedAt = t
		art.Source = a.Name()
		*as = append(*as, art)
	})

	return nil
}
