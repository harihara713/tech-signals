package source

import (
	"fmt"
	"strings"
	"sync"
	"unicode"

	"github.com/harihara713/tech-signals/article"
)

type Source interface {
	Name() string
	Fetch(as *article.ArticleStore, wg *sync.WaitGroup) error
}

func transformToTimeString(str string) (string, error) {
	parts := strings.Fields(str)
	cleaned := make([]string, 0, len(parts))
	for _, w := range parts {
		cleanPart := strings.TrimFunc(w, unicode.IsPunct)

		if cleanPart != "" {
			cleaned = append(cleaned, cleanPart)
		}
	}

	if len(cleaned) != 3 {
		return "", fmt.Errorf("some parts missing in date string")
	}

	// june or JUNE to June
	// i am assuming the input have month first with format like Sept, june, oct like this
	// next part is date in numeric form, and last part is year
	month := cleaned[0]
	date := cleaned[1]
	year := cleaned[2]

	month = strings.ToUpper(string(month[0])) + strings.ToLower(month[1:])
	if len(date) == 1 {
		date = "0" + date
	}

	return fmt.Sprintf("%s %s, %s", month, date, year), nil
}
