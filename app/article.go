package app

import (
	"time"
)

// Article represents a single article
type Article struct {
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
