package source

import "github.com/harihara713/tech-signals/article"

type Source interface {
	Name() string
	Fetch(as *article.ArticleStore) error
}
