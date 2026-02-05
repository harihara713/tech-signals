package source

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/harihara713/tech-signals/article"
)

type Github struct {
}

func (g Github) Name() string {
	return "Github"
}

func (g Github) Fetch(as *article.ArticleStore, wg *sync.WaitGroup) error {
	url := "https://github.blog/engineering/"
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Github: Error getting response", "err", err)
		wg.Done()
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Github: Error response", "status", resp.StatusCode)
		wg.Done()
		return fmt.Errorf("Github blogs response status: '%d'", resp.StatusCode)
	}
	// extract the articles details
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		slog.Error("Github: Error parsing the page", "err", err)
		wg.Done()
		return err
	}

	doc.Find(".post-columns.post-columns--3-3 article").Each(func(i int, s *goquery.Selection) {
		art := article.Article{}
		div := s.Find("div")
		secondChild := div.ChildrenFiltered("div").Eq(1)
		link := secondChild.Find("h3 a")

		art.URL = strings.TrimSpace(link.AttrOr("href", ""))
		art.Title = strings.TrimSpace(link.Text())
		art.Summary = secondChild.Find("div p").Text()

		footer := s.Find("footer div span")

		art.Author = footer.Find("span a").Text()
		timeStr := strings.TrimSpace(footer.Find("time").Text())

		t, err := time.Parse("January 02, 2006", timeStr)
		if err != nil {
			slog.Error("Github: Error Parsing time", "err", err)
			return
		}

		art.PublishedAt = t
		art.Source = g.Name()

		*as = append(*as, art)
	})

	wg.Done()
	return nil
}
