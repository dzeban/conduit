package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

// List returns n articles from Postgres
func (s Store) List(f app.ArticleListFilter) ([]app.Article, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err :=
		psql.Select(`
				a.slug as slug,
				a.title as title,
				a.description as description,
				a.body as body,
				a.created as created,
				a.updated as updated,
				a.author as username,
				u.bio as bio,
				u.image as image,
				f.follows != '' as following
			`).
			From("articles a").
			Join("users u on (a.author=u.name)").
			LeftJoin("followers f on (u.name=f.follows)").
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

	var title, slug, authorName string
	var description, body, bio, image sql.NullString
	var created, updated time.Time
	var following sql.NullBool
	for rows.Next() {
		err = rows.Scan(
			&slug, &title, &description, &body, &created, &updated,
			&authorName, &bio, &image, &following,
		)
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
			Author: app.Profile{
				Name:      authorName,
				Bio:       bio.String,
				Image:     image.String,
				Following: following.Bool,
			},
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (s Store) Feed(f app.ArticleListFilter) ([]app.Article, error) {
	// Articles feed query
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err :=
		psql.Select(`
				a.slug as slug,
				a.title as title,
				a.description as description,
				a.body as body,
				a.created as created,
				a.updated as updated,
				a.author as username,
				u.bio as bio,
				u.image as image,
				f.follows != '' as following
			`).
			From("articles a").
			Join("users u on (a.author=u.name)").
			LeftJoin("followers f on (u.name=f.follows)").
			Where("f.follower=?", f.Username).
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

	var title, slug, authorName string
	var description, body, bio, image sql.NullString
	var created, updated time.Time
	var following sql.NullBool
	for rows.Next() {
		err = rows.Scan(
			&slug, &title, &description, &body, &created, &updated,
			&authorName, &bio, &image, &following,
		)
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
			Author: app.Profile{
				Name:      authorName,
				Bio:       bio.String,
				Image:     image.String,
				Following: following.Bool,
			},
		}

		articles = append(articles, article)
	}

	return articles, nil
}

// Get returns a single article by its slug
func (s Store) Get(slug string) (*app.Article, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err :=
		psql.Select(`
				a.title as title,
				a.description as description,
				a.body as body,
				a.created as created,
				a.updated as updated,
				a.author as username,
				u.bio as bio,
				u.image as image,
				f.follows != '' as following
			`).
			From("articles a").
			Join("users u on (a.author=u.name)").
			LeftJoin("followers f on (u.name=f.follows)").
			Where(sq.Eq{"slug": slug}).
			ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build select query")
	}

	row := s.db.QueryRowx(query, args...)

	var title, authorName string
	var description, body, bio, image sql.NullString
	var created, updated time.Time
	var following sql.NullBool

	err = row.Scan(
		&title, &description, &body, &created, &updated,
		&authorName, &bio, &image, &following,
	)
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
		Author: app.Profile{
			Name:      authorName,
			Bio:       bio.String,
			Image:     image.String,
			Following: following.Bool,
		},
	}

	return &article, nil
}

func (s Store) Create(a *app.Article) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err :=
		psql.
			Insert("articles").
			Columns("slug", "title", "description", "body", "author").
			Values(a.Slug, a.Title, a.Description, a.Body, a.Author.Name).
			ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build insert query")
	}

	fmt.Println(query, args)
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute insert query")
	}

	return nil
}

func (s Store) Delete(slug string) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err :=
		psql.
			Delete("articles").
			Where(sq.Eq{"slug": slug}).
			ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build delete query")
	}

	fmt.Println(query, args)
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute delete query")
	}

	return nil
}

func (s Store) Update(slug string, a *app.Article) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err :=
		psql.
			Update("articles").
			SetMap(map[string]interface{}{
				"title":       a.Title,
				"description": a.Description,
				"body":        a.Body,
			}).
			Where(sq.Eq{"slug": slug}).
			ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build update query")
	}

	fmt.Println(query, args)
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute update query")
	}

	return nil
}
