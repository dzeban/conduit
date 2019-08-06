package article

import (
	"database/sql"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/db"
)

// Service implements app.ArticleService interface
// It serves articles from Postgres
type Service struct {
	db     *sqlx.DB
	router *chi.Mux
}

// New creates new Article service backed by Postgres
func NewService(DSN string) (*Service, error) {
	db, err := db.ConnectLoop("postgres", DSN, 1*time.Minute)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to articles db")
	}

	router := chi.NewRouter()

	s := &Service{db: db, router: router}
	router.Get("/articles/", s.HandleArticles)
	router.Get("/articles/{slug}", s.HandleArticle)

	return s, nil
}

// List returns n articles from Postgres
func (s *Service) List(n int) ([]app.Article, error) {
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
func (s *Service) Get(slug string) (*app.Article, error) {
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
