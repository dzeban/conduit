package mock

import (
	"errors"
	"time"

	"github.com/dzeban/conduit/app"
)

var (
	ArticleValid = app.Article{
		Id:          1,
		Title:       "Title",
		Description: "Description",
		Body:        "Body",
		Author:      Author,
		Created:     time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		Updated:     time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
	}

	ArticleUpdated = app.Article{
		Id:          2,
		Title:       "Other title",
		Description: "Other description",
		Body:        "Other body",
		Author:      Author,
		Created:     time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
		Updated:     time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
	}

	Author = app.Profile{
		Id:   1,
		Name: "author",
		Bio:  "bio",
	}
)

// ArticleStore is a fake implementation of article.Store as Go map
type ArticleStore struct {
	ById   map[int]*app.Article
	BySlug map[string]*app.Article
}

func NewArticleStore() *ArticleStore {
	as := &ArticleStore{
		ById:   make(map[int]*app.Article),
		BySlug: make(map[string]*app.Article),
	}

	_ = as.CreateArticle(&ArticleValid)
	_ = as.CreateArticle(&ArticleUpdated)

	return as
}

func (as *ArticleStore) CreateArticle(a *app.Article) error {
	as.ById[a.Id] = a
	as.BySlug[a.Slug] = a
	return nil
}

func (as *ArticleStore) ListArticles(f *app.ArticleListFilter) ([]*app.Article, error) {
	panic("not implemented") // TODO: Implement
}

func (as *ArticleStore) GetArticle(slug string) (*app.Article, error) {
	return as.BySlug[slug], nil
}

func (as *ArticleStore) UpdateArticle(a *app.Article) error {
	as.ById[a.Id] = a
	as.BySlug[a.Slug] = a
	return nil
}

func (as *ArticleStore) DeleteArticle(id int) error {
	a, ok := as.ById[id]
	if !ok {
		return errors.New("not found by id")
	}

	delete(as.ById, id)
	delete(as.BySlug, a.Slug)

	return nil
}
