package article

import (
	"database/sql"
	"time"

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
func (s PostgresStore) List(n int) ([]app.Article, error) {
	queryArticles := `
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
		LIMIT $1
	`

	rows, err := s.db.Queryx(queryArticles, n)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query articles")
	}

	var articles []app.Article
	for rows.Next() {
		var article app.Article
		err = rows.StructScan(&article)
		if err != nil {
			// TODO: log.Println(err)
			continue
		}

		articles = append(articles, article)
	}

	return articles, nil
}

// Get returns a single article by its slug
func (s PostgresStore) Get(slug string) (*app.Article, error) {
	queryArticle := `
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

	row := s.db.QueryRowx(queryArticle, slug)

	var article app.Article
	err := row.StructScan(&article)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to query article")
	}

	return &article, nil
}
