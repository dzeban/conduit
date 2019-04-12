package postgres

import (
	"database/sql"

	"github.com/dzeban/conduit/app"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// ArticlesService implements app.ArticleService interface
// It serves articles from Postgres
type ArticlesService struct {
	db *sqlx.DB
}

// New creates new Articles service backed by Postgres
func New(DSN string) (*ArticlesService, error) {
	db, err := sqlx.Connect("postgres", DSN)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to articles db")
	}

	return &ArticlesService{db: db}, nil
}

// List returns n articles from Postgres
func (s *ArticlesService) List(n int) ([]app.Article, error) {
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
func (s *ArticlesService) Get(slug string) (*app.Article, error) {
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
