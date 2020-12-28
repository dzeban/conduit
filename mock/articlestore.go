package mock

import (
	"time"

	"github.com/dzeban/conduit/app"
)

var (
	ArticleValid = app.Article{
		Title:       "Title",
		Description: "Description",
		Body:        "Body",
		Author:      Author,
		Created:     time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		Updated:     time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
	}

	ArticleUpdated = app.Article{
		Title:       "Other title",
		Description: "Other description",
		Body:        "Other body",
		Author:      Author,
		Created:     time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
		Updated:     time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
	}

	Author = app.Profile{
		Name: "author",
		Bio:  "bio",
	}
)

// ArticleStore is a fake implementation of article.Store as Go map
type ArticleStore struct {
	m map[string]*app.Article
}

func NewArticleStore() *ArticleStore {
	as := &ArticleStore{
		m: make(map[string]*app.Article),
	}

	_ = as.CreateArticle(&ArticleValid)
	_ = as.CreateArticle(&ArticleUpdated)

	return as
}

func (as *ArticleStore) CreateArticle(a *app.Article) error {
	as.m[a.Slug] = a
	return nil
}

func (as *ArticleStore) ListArticles(f app.ArticleListFilter) ([]app.Article, error) {
	panic("not implemented") // TODO: Implement
}

func (as *ArticleStore) GetArticle(slug string) (*app.Article, error) {
	return as.m[slug], nil
}

func (as *ArticleStore) UpdateArticle(slug string, a *app.Article) error {
	as.m[slug] = a
	return nil
}

func (as *ArticleStore) DeleteArticle(slug string) error {
	panic("not implemented") // TODO: Implement
}
