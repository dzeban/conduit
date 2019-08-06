package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/koding/multiconfig"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/service/article"
	"github.com/dzeban/conduit/service/user"
)

type ServerConfig struct {
	Port int `default:"8080"`
}

// Config represents app configuration
type Config struct {
	Server   ServerConfig
	Articles app.ArticleServiceConfig
	Users    app.UserServiceConfig
}

func main() {
	var config Config
	multiconfig.New().MustLoad(&config)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Server.Port),
		Handler: router,
	}

	log.Printf("using config: %#v\n", config)

	articleService, err := article.NewService(config.Articles.DSN)
	if err != nil {
		log.Fatal("cannot create articles service: ", err)
	}

	userService, err := user.NewFromDSN(config.Users.DSN, config.Users.Secret)
	if err != nil {
		log.Fatal("cannot create users service: ", err)
	}

	// Setup API endpoints
	router.Mount("/articles", articleService)
	router.Mount("/users", userService)
	router.Mount("/profiles", userService)

	log.Println("start listening on", server.Addr)
	log.Fatal(server.ListenAndServe())
}
