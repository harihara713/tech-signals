package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/harihara713/tech-signals/article"
	"github.com/harihara713/tech-signals/source"
)

var sourceNames = []string{"meta", "google", "amazon", "github", "stripe"}

func main() {
	var (
		selectedSource string
		limit          int
		saveFilename   string
	)

	flag.StringVar(&selectedSource, "source", "", "source name (Meta, Google, Amazon, Github, Stripe)")
	flag.IntVar(&limit, "limit", 0, "max number of blogs to fetch (0 = all)")
	flag.StringVar(&saveFilename, "save", "", "output filename (json)")

	flag.Parse()

	if selectedSource != "" {
		selectedSource = strings.ToLower(selectedSource)
		var validSource bool
		if slices.Contains(sourceNames, selectedSource) {
			validSource = true
		}

		if !validSource {
			fmt.Fprintf(os.Stderr, "%s: invalid source name,\nPlease provide valid source name(meta, google, amazon, github, stripe)", selectedSource)
			os.Exit(1)
		}

		var articles *article.ArticleStore
		// fetch all the blogs for a particular source
		switch selectedSource {
		case "meta":
			source.Meta{}.Fetch(articles)
		case "google":
			source.Google{}.Fetch(articles)
		case "amazon":
			source.Amazon{}.Fetch(articles)
		case "github":
			source.Github{}.Fetch(articles)
		case "stripe":
			source.Stripe{}.Fetch(articles)
		}

		fmt.Printf("\nTotal %d articles fetched\n\n", len(*articles))
		for _, art := range *articles {
			fmt.Printf("[%s]\nTitle: %s\nUrl: %s\nPublished At: %v\nAuthor: %s\nSummary: %s\nTags: %s\n\n", art.Source, art.Title,
				art.URL, art.PublishedAt, art.Author, art.Summary, strings.Join(art.Tags, ", "))
		}
	}

}
