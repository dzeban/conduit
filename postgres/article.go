package postgres

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

// PostgresArticle is the same as app.Article but specific for articles list
// queries. It expands Author type to its fields and uses sql.Null* types.
// All of this is needed to allow scanning query result.
type PostgresArticle struct {
	Id          int
	Slug        string
	Title       string
	Description sql.NullString
	Body        sql.NullString
	Created     time.Time
	Updated     time.Time
	AuthorId    int
	AuthorName  string
	AuthorBio   string
	AuthorImage sql.NullString
	Following   bool
}

func (s Store) ListArticles(f *app.ArticleListFilter) ([]*app.Article, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q := psql.Select(`
				a.id as id,
				a.slug as slug,
				a.title as title,
				a.description as description,
				a.body as body,
				a.created as created,
				a.updated as updated,
				a.author_id as author_id,
				u.name as author_name,
				u.bio as author_bio,
				u.image as author_image
			`).
		From("articles a").
		Join("users u on (a.author_id = u.id)").
		OrderBy("created DESC")

	if f.CurrentUser != nil {
		q = q.LeftJoin("followers f on (u.id = f.followee)").
			Where("f.follower = ?", f.CurrentUser.Id).
			Columns("f.followee != 0 as following")
	}

	if f.Author != nil {
		q = q.Where("author_id = ?", f.Author.Id)
	}
	// q = q.Where("favorite = ?", f.Favorite)

	q = q.Limit(f.Limit).Offset(f.Offset)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build select query")
	}

	rows, err := s.db.Queryx(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query articles")
	}
	defer rows.Close()

	// We need to scan each row to the struct. But there are 2 problems here:
	//
	// 1. Some fields may be Null and thus need to be handled with sql.Null*
	//    types. Structs that we scan into are from application level and don't
	//    use sql specific types.
	// 2. Columns set is dynamic. When query is performed with followers table
	// 	  the result set has the "following" column.
	//
	// Because of this we use special PostgresArticle type that allows scanning
	// directly into the struct and then create app.Article values from it.
	var (
		a        PostgresArticle
		articles []*app.Article
	)
	for rows.Next() {
		err := rows.StructScan(&a)
		if err != nil {
			return nil, errors.Wrap(err, "row scan failed")
		}

		articles = append(articles, &app.Article{
			Id:          a.Id,
			Slug:        a.Slug,
			Title:       a.Title,
			Description: a.Description.String,
			Body:        a.Body.String,
			Created:     a.Created,
			Updated:     a.Updated,
			Author: app.Profile{
				Id:    a.AuthorId,
				Name:  a.AuthorName,
				Image: a.AuthorImage.String,
			},
		})
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows scan error")
	}

	return articles, nil
}

// Get returns a single article by its slug
func (s Store) GetArticle(slug string) (*app.Article, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err :=
		psql.Select(`
				a.title as title,
				a.description as description,
				a.body as body,
				a.created as created,
				a.updated as updated,
				a.author_id as author_id,
				u.name as author_name,
				u.bio as bio,
				u.image as image,
				f.followee != 0 as following
			`).
			From("articles a").
			Join("users u on (a.author_id = u.id)").
			LeftJoin("followers f on (u.id = f.followee)").
			Where(sq.Eq{"a.slug": slug}).
			ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build select query")
	}

	row := s.db.QueryRowx(query, args...)

	// TODO: use PostgresArticle with sqlx.StructScan
	var title, authorName string
	var authorId int
	var description, body, bio, image sql.NullString
	var created, updated time.Time
	var following sql.NullBool

	err = row.Scan(
		&title, &description, &body, &created, &updated,
		&authorId, &authorName, &bio, &image, &following,
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
			Id:        authorId,
			Name:      authorName,
			Bio:       bio.String,
			Image:     image.String,
			Following: following.Bool,
		},
	}

	return &article, nil
}

func (s Store) CreateArticle(a *app.Article) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err :=
		psql.
			Insert("articles").
			Columns("slug", "title", "description", "body", "author_id", "created", "updated").
			Values(a.Slug, a.Title, a.Description, a.Body, a.Author.Id, a.Created, a.Updated).
			ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build insert query")
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute insert query")
	}

	return nil
}

func (s Store) DeleteArticle(id int) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err :=
		psql.
			Delete("articles").
			Where(sq.Eq{"id": id}).
			ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build delete query")
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute delete query")
	}

	return nil
}

func (s Store) UpdateArticle(a *app.Article) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err :=
		psql.
			Update("articles").
			SetMap(a.UpdateMap()).
			Where(sq.Eq{"id": a.Id}).
			ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build update query")
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute update query")
	}

	return nil
}
