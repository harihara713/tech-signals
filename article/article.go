package article

import "time"

type Article struct {
	Title       string
	URL         string
	Summary     string
	Author      string
	Tags        []string
	PublishedAt time.Time
	Source      string
}

type ArticleStore []Article

func NewArticleStore() *ArticleStore {
	return &ArticleStore{}
}
