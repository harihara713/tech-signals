package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/harihara713/tech-signals/article"
	"github.com/harihara713/tech-signals/source"
)

func main() {
	sources := []source.Source{
		source.Meta{},
	}

	artStore := article.NewArticleStore()

	for _, s := range sources {
		if err := s.Fetch(artStore); err != nil {
			slog.Error("Error fetching articles", "err", err)
		}
	}

	// show
	fmt.Println("\n" + strings.Repeat("=", 100) + "\n")
	fmt.Printf("Total %d articles found\n", len(*artStore))

	for _, a := range *artStore {
		fmt.Printf("%#v\n", a)
	}
}
