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
}

// ArticlesService defines an interface to work with articles
type ArticlesService interface {
	List(n int) ([]Article, error)
	Get(slug string) (*Article, error)
}

// ArticlesConfig describes configuration for ArticleService
type ArticlesConfig struct {
	Type string `default:"postgres"`
	DSN  string `default:"postgres://postgres:postgres@localhost/conduit?sslmode=disable"`
}
