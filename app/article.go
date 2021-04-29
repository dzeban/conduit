package app

import (
	"time"

	"github.com/pkg/errors"
)

// Article represents a single article
type Article struct {
	Id          int       `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Author      Profile   `json:"author"`

	// TagList []Tag `json:"tagList"`
	// IsFavorited bool `json:"favorited"`
	// FavoritesCount int `json:"favoritesCount"`
}

// UpdateMap returns map of fields to be updated. It sets only subset of fields
// that are allowed to be updated and set it if it contains non empty values.
func (a Article) UpdateMap() map[string]interface{} {
	m := make(map[string]interface{})

	if a.Title != "" {
		m["title"] = a.Title
	}
	if a.Description != "" {
		m["description"] = a.Description
	}
	if a.Body != "" {
		m["body"] = a.Body
	}

	emptyTime := time.Time{}
	if a.Updated != emptyTime {
		m["updated"] = a.Updated
	}

	return m
}

// ArticleServiceConfig describes configuration for ArticleService
type ArticleServiceConfig struct {
	Type   string `default:"postgres"`
	DSN    string `default:"postgres://postgres:postgres@postgres/conduit?sslmode=disable"`
	Secret string
}

type ArticleListFilter struct {
	CurrentUser *User // used for favorites and following filtering
	Author      *Profile
	Limit       uint64
	Offset      uint64
}

// NewArticleListFilter creates filter with default values
func NewArticleListFilter() ArticleListFilter {
	return ArticleListFilter{
		Limit:  20,
		Offset: 0,
	}
}

func (f ArticleListFilter) Validate() error {
	// Silly filters to save database from huge queries
	if f.Limit > 100 || f.Offset > 10000 {
		return errors.New("invalid article list filter")
	}
	return nil
}
