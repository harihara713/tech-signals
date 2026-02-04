package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

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

	articles := article.NewArticleStore()
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

		// fetch all the blogs for a particular source
		switch selectedSource {
		case "meta":
			fetchArticles(articles, source.Meta{})
		case "google":
			fetchArticles(articles, source.Google{})
		case "amazon":
			fetchArticles(articles, source.Amazon{})
		case "github":
			fetchArticles(articles, source.Github{})
		case "stripe":
			fetchArticles(articles, source.Stripe{})
		}

		if articles == nil {
			fmt.Println("\nTotal 0 articles fetched")
			fmt.Println()
			return
		}

		fmt.Printf("\nTotal %d articles fetched\n\n", len(*articles))
		for _, art := range *articles {
			fmt.Printf("[%s]\nTitle: %s\nUrl: %s\nPublished At: %v\nAuthor: %s\nSummary: %s\nTags: %s\n\n", art.Source, art.Title,
				art.URL, art.PublishedAt, art.Author, art.Summary, strings.Join(art.Tags, ", "))
		}

		return
	}

	sources := []source.Source{
		source.Amazon{}, source.Github{}, source.Google{}, source.Meta{}, source.Stripe{},
	}

	fetchArticles(articles, sources...)

	if limit > 0 {
		if limit > len(*articles) {
			limit = len(*articles)
		}

		limitedResult := (*articles)[:limit]

		for _, art := range limitedResult {
			fmt.Printf("[%s]\nTitle: %s\nUrl: %s\nPublished At: %v\nAuthor: %s\nSummary: %s\nTags: %s\n\n", art.Source, art.Title,
				art.URL, art.PublishedAt, art.Author, art.Summary, strings.Join(art.Tags, ", "))
		}

		return
	}

	if saveFilename != "" {
		saveFilename = filepath.Base(saveFilename)
		data, err := json.MarshalIndent(articles, "", " ")
		if err != nil {
			fmt.Println("Error: Failed to marshal the articles")
			return
		}

		f, err := os.OpenFile(saveFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("Error: Failed to create the file")
			return
		}
		defer f.Close()
		if _, err := f.Write(data); err != nil {
			fmt.Println("Error: Failed to write to the file")
			return
		}

		fmt.Println("Successfully append to the file")
		return
	}

	fmt.Printf("\nTotal %d articles fetched\n\n", len(*articles))
	for _, art := range *articles {
		fmt.Printf("[%s]\nTitle: %s\nUrl: %s\nPublished At: %v\nAuthor: %s\nSummary: %s\nTags: %s\n\n", art.Source, art.Title,
			art.URL, art.PublishedAt, art.Author, art.Summary, strings.Join(art.Tags, ", "))
	}
}

func fetchArticles(articles *article.ArticleStore, sources ...source.Source) {
	// Check the internet connection
	_, err := net.DialTimeout("tcp", "google.com:80", 3*time.Second)
	if err != nil {
		fmt.Printf("Error: Please check your connection")
		return
	}
	// fetch the articles

	for _, s := range sources {
		s.Fetch(articles)
	}
}
