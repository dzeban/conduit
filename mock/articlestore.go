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

	ArticleToDelete = app.Article{
		Id:          4,
		Title:       "Deleted title",
		Slug:        "deleted-title-azxs",
		Description: "Deleted description",
		Body:        "Deleted body",
		Author:      Profile1,
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
	_ = as.CreateArticle(&Article3)
	_ = as.CreateArticle(&ArticleToDelete)

	return as
}

func (as *ArticleStore) CreateArticle(a *app.Article) error {
	as.ById[a.Id] = a
	as.BySlug[a.Slug] = a
	return nil
}

func (as *ArticleStore) ListArticles(f *app.ArticleListFilter) ([]*app.Article, error) {
	articles := []*app.Article{&ArticleValid, &ArticleUpdated, &Article3}

	// Mock feed for one user
	if f.CurrentUser != nil && f.CurrentUser.Id == UserValid.Id {
		articles = []*app.Article{&ArticleValid, &ArticleUpdated}
	}

	if f.Author != nil {
		switch f.Author.Id {
		case Profile1.Id:
			articles = []*app.Article{&ArticleValid, &ArticleUpdated}

		case Profile2.Id:
			articles = []*app.Article{&Article3}

		default:
			return nil, nil
		}
	}

	var result []*app.Article
	if f.Offset > len(articles)-1 {
		return nil, nil
	}

	limit := f.Limit
	if limit > len(articles) {
		limit = len(articles)
	}
	for i := f.Offset; i < limit; i++ {
		result = append(result, articles[i])
	}

	return result, nil
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
