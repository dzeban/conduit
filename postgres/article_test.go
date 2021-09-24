package postgres

import (
	"os"
	"testing"
	"time"

	"github.com/dzeban/conduit/app"
)

// shouldSkip check environment variable CONDUIT_TEST_DB and if it's not set it
// will skip the test identified by state t.
// This is needed to avoid expensive system tests by default without having
// build tags.
func shouldSkip(t *testing.T) {
	_, ok := os.LookupEnv("CONDUIT_TEST_DB")
	if !ok {
		t.Skip("CONDUIT_TEST_DB not set, skipping system tests")
	}
}

func TestArticleStoreSimple(t *testing.T) {
	shouldSkip(t)

	s, err := NewStore("postgres://conduit:conduit@localhost:5432/conduit?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	a1 := app.Article{
		Slug:        "test-article-store-1-slug",
		Title:       "Test article store 1 title",
		Description: "Test article store 1 description",
		Body:        "Test article store 1 body",
		Created:     time.Now(),
		Updated:     time.Now(),
		Author: app.Profile{
			Id: 1,
		},
	}

	defer func() {
		_ = s.DeleteArticleBySlug(a1.Slug)
	}()

	err = s.CreateArticle(&a1)
	if err != nil {
		t.Fatalf("s.CreateArticle(%v) failed: %v", a1, err)
	}
	a, err := s.GetArticle(a1.Slug)
	if err != nil {
		t.Fatalf("s.GetArticle(%v) failed: %v", a1.Slug, err)
	}

	err = s.UpdateArticle(a)
	if err != nil {
		t.Fatalf("s.UpdateArticle(%v) failed: %v", a, err)
	}

	err = s.DeleteArticle(a.Id)
	if err != nil {
		t.Fatalf("s.DeleteArticle(%v) failed: %v", a.Id, err)
	}
}
