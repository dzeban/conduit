package article

import (
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/db"
)

type PostgresStore struct {
	db *sqlx.DB
}

func New(DSN string) (app.ArticleStore, error) {
	db, err := db.ConnectLoop("postgres", DSN, 1*time.Minute)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to users db")
	}

	return PostgresStore{db: db}, nil
}

// List returns n articles from Postgres
func (s PostgresStore) List(f app.ArticleListFilter) ([]app.Article, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err :=
		psql.Select("slug, title, description, body, created, updated").
			From("articles").
			Where(f.Map()).
			Limit(f.Limit).
			Offset(f.Offset).
			ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build select query")
	}

	rows, err := s.db.Queryx(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query articles")
	}

	var articles []app.Article

	var title, slug string
	var description, body sql.NullString
	var created, updated time.Time
	for rows.Next() {
		err = rows.Scan(&slug, &title, &description, &body, &created, &updated)
		if err != nil {
			log.Println(err)
			continue
		}

		article := app.Article{
			Slug:        slug,
			Title:       title,
			Description: description.String,
			Body:        body.String,
			Created:     created,
			Updated:     updated,
		}

		articles = append(articles, article)
	}

	return articles, nil
}

// Get returns a single article by its slug
func (s PostgresStore) Get(slug string) (*app.Article, error) {
	queryArticle := `
		SELECT
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

	row := s.db.QueryRowx(queryArticle, slug)

	var title string
	var description, body sql.NullString
	var created, updated time.Time

	err := row.Scan(&title, &description, &body, &created, &updated)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to query article")
	}

	article := app.Article{
		Slug:        slug,
		Title:       title,
		Description: description.String,
		Body:        body.String,
		Created:     created,
		Updated:     updated,
	}

	return &article, nil
}
