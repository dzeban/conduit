package main

import (
	"time"

	_ "github.com/lib/pq"
)

const (
	queryArticle = `
	SELECT
		slug,
		title,
		description,
		body,
		created,
		updated
	FROM
		articles
	WHERE
		slug = $1
	`

	queryArticles = `
	SELECT
		slug,
		title,
		description,
		body,
		created,
		updated
	FROM
		articles
	ORDER BY
		created DESC
	LIMIT
		20
	`
)

// Article represents a single article
type Article struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Slug        string    `json:"slug,omitempty"`
	Body        string    `json:"body,omitempty"`
	Created     time.Time `json:"created,omitempty"`
	Updated     time.Time `json:"updated,omitempty"`
}

// Articles is a dedicated type for a list of articles
type Articles []Article
