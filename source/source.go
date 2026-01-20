package source

import "github.com/harihara713/goscrape/article"

type Source interface {
	Name() string
	Fetch(as *article.ArticleStore) error
}
