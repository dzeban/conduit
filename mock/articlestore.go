package mock

import (
	"errors"
	"time"

	"github.com/dzeban/conduit/app"
)

var (
	ArticleValid = app.Article{
		Id:          1,
		Slug:        "title-q1w2",
		Title:       "Title",
		Description: "Description",
		Body:        "Body",
		Author:      Profile1,
		Created:     time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		Updated:     time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
	}

	ArticleUpdated = app.Article{
		Id:          2,
		Title:       "Other title",
		Slug:        "other-title-azxs",
		Description: "Other description",
		Body:        "Other body",
		Author:      Profile1,
		Created:     time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
		Updated:     time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
	}

	Article3 = app.Article{
		Id:          3,
		Title:       "Another title",
		Slug:        "another-title-qwer",
		Description: "Another description",
		Body:        "Another body",
		Author:      Profile2,
		Created:     time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
		Updated:     time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
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
	if f.Author != nil {
		switch f.Author.Id {
		case Profile1.Id:
			return []*app.Article{&ArticleValid, &ArticleUpdated}, nil

		case Profile2.Id:
			return []*app.Article{&Article3}, nil

		default:
			return nil, nil
		}
	}

	return []*app.Article{&ArticleValid, &ArticleUpdated, &Article3}, nil
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
