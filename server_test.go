package main

import (
	"testing"

	"github.com/dzeban/conduit/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
