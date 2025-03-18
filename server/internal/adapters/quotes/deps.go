package quotes

import (
	"math/rand"

	svc "github.com/klynxe/word-of-wisdom/server/internal/services/quotes"
	quotesStorage "github.com/klynxe/word-of-wisdom/server/internal/storage/quotes"
)

type Deps struct {
	quotes []string
}

var _ svc.Deps = (*Deps)(nil)

func NewDeps(storage *quotesStorage.Storage) *Deps {
	return &Deps{
		quotes: storage.GetAllQuotes(),
	}
}

func (d *Deps) GetRandomQuote() string {
	return d.quotes[rand.Intn(len(d.quotes))]
}
