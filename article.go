package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
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

func articleHandle(c *gin.Context) {
	connStr := "postgres://postgres:postgres@localhost/conduit?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic("can't connect to db")
	}

	var article Article
	row := db.QueryRow(queryArticle, c.Param("slug"))
	err = row.Scan(&article.Slug, &article.Title, &article.Description, &article.Body, &article.Created, &article.Updated)
	if err == sql.ErrNoRows {
		c.Status(404)
		return
	} else if err != nil {
		c.Status(500)
		return
	}

	c.JSON(200, article)
}

func articlesHandle(c *gin.Context) {
	connStr := "postgres://postgres:postgres@localhost/conduit?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic("can't connect to db")
	}

	rows, err := db.Query(queryArticles)
	if err != nil {
		c.Status(500)
		return
	}

	var articles Articles
	for rows.Next() {
		var article Article
		err = rows.Scan(&article.Slug, &article.Title, &article.Description, &article.Body, &article.Created, &article.Updated)
		if err != nil {
			log.Println(err)
			continue
		}

		articles = append(articles, article)
	}

	c.JSON(200, articles)
}
