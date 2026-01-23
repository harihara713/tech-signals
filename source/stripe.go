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

type Stripe struct {
}

func (s Stripe) Name() string {
	return "Stripe"
}

func (s Stripe) Fetch(as *article.ArticleStore) error {
	url := "https://stripe.com/blog/engineering"
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Stripe: Error getting response", "err", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Stripe: Error response", "status", resp.StatusCode)
		return fmt.Errorf("Stripe blogs response status: '%d'", resp.StatusCode)
	}
	// extract the articles details
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		slog.Error("Stripe: Error parsing the page", "err", err)
		return err
	}

	doc.Find(".BlogIndexGridSection__layout article").Each(func(i int, sl *goquery.Selection) {
		art := article.Article{}

		art.Title = strings.TrimSpace(sl.Find("h1").Text())
		link := sl.Find("time a")
		art.URL = strings.TrimSpace(link.AttrOr("href", ""))
		timeStr := strings.TrimSpace(link.Text())

		t, err := time.Parse("January 02, 2006", timeStr)
		if err != nil {
			slog.Error("Stripe: Error Parsing time", "err", err)
			return
		}

		art.Author = sl.Find(".BlogIndexPost__authorList figcaption a").Text()
		art.Summary = sl.Find(".BlogIndexPost__body p").Text()

		art.PublishedAt = t
		art.Source = s.Name()

		*as = append(*as, art)
	})

	return nil
}
