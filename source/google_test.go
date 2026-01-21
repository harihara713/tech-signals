package source

import (
	"testing"

	"github.com/harihara713/tech-signals/article"
)

func TestGoogle(t *testing.T) {
	g := Google{}

	a := article.NewArticleStore()
	g.Fetch(a)
}
