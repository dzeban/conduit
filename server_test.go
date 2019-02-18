package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/dzeban/conduit/mock"
)

func TestArticlesList(t *testing.T) {
	mockArticlesService, _ := mock.New()
	s := Server{
		articles: mockArticlesService,
	}

	router := gin.New()
	router.GET("/articles/", s.HandleArticles)

	assert.HTTPRedirect(t, router.ServeHTTP, "GET", "/articles", nil)
	assert.HTTPSuccess(t, router.ServeHTTP, "GET", "/articles/", nil)
	assert.HTTPBodyContains(t, router.ServeHTTP, "GET", "/articles/", nil, "Title 1")
}

func TestNotFoundArticle(t *testing.T) {
	mockArticlesService, _ := mock.New()
	s := Server{
		articles: mockArticlesService,
	}

	router := gin.New()
	router.GET("/articles/", s.HandleArticles)

	assert.HTTPError(t, router.ServeHTTP, "GET", "/articles/xxx", nil)
}
