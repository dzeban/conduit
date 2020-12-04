package app

import (
	"time"
)

// Article represents a single article
type Article struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Slug        string    `json:"slug"`
	Body        string    `json:"body"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Author      Profile   `json:"author"`
}

type ArticleCreateRequest struct {
	Article Article `json:"article"`
}

type ArticleUpdateRequest struct {
	Article ArticleUpdateRequestData `json:"article"`
}

type ArticleUpdateRequestData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

// ArticleStore defines an interface to work with articles
type ArticleStore interface {
	List(f ArticleListFilter) ([]Article, error)
	Feed(f ArticleListFilter) ([]Article, error)
	Get(slug string) (*Article, error)
	Create(a *Article) error
	Update(a *Article) error
	Delete(slug string) error
}

// ArticleServiceConfig describes configuration for ArticleService
type ArticleServiceConfig struct {
	Type   string `default:"postgres"`
	DSN    string `default:"postgres://postgres:postgres@postgres/conduit?sslmode=disable"`
	Secret string
}

type ArticleListFilter struct {
	Username string
	Limit    uint64
	Offset   uint64
}

// NewArticleListFilter creates filter with default values
func NewArticleListFilter() ArticleListFilter {
	return ArticleListFilter{
		Limit:  20,
		Offset: 0,
	}
}

func (f ArticleListFilter) Map() map[string]interface{} {
	m := make(map[string]interface{})
	if f.Username != "" {
		m["username"] = f.Username
	}

	return m
}
