package article

import (
	"time"

	"github.com/dzeban/conduit/app"
)

type mockStore struct {
	articles map[string]app.Article
}

func newMockStore() app.ArticleStore {
	article := app.Article{
		Slug:    "test",
		Title:   "Test title",
		Body:    "Test body of the article. Several sentences.",
		Created: time.Now(),
		Updated: time.Now(),
	}
	articles := make(map[string]app.Article)
	articles[article.Slug] = article

	return mockStore{articles}
}

func (s mockStore) Get(slug string) (*app.Article, error) {
	a := s.articles[slug] // because go complains: cannot take the address of s.articles[slug]
	return &a, nil
}

func (s mockStore) List(f app.ArticleListFilter) ([]app.Article, error) {
	i := uint64(0)
	var list []app.Article
	for _, v := range s.articles {
		list = append(list, v)
		i++
		if i >= f.Limit {
			break
		}
	}

	return list, nil
}

func (s mockStore) Feed(f app.ArticleListFilter) ([]app.Article, error) {
	i := uint64(0)
	var list []app.Article
	for _, v := range s.articles {
		list = append(list, v)
		i++
		if i >= f.Limit {
			break
		}
	}

	return list, nil
}
