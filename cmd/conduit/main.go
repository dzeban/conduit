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
	"github.com/dzeban/conduit/postgres"
	"github.com/dzeban/conduit/user"
)

// Config represents app configuration
type Config struct {
	Server   ServerConfig
	Articles app.ArticleServiceConfig
	Users    UserServiceConfig
}

// UserServiceConfig describes configuration for UserService
type UserServiceConfig struct {
	DSN    string `default:"postgres://postgres:postgres@postgres/conduit?sslmode=disable"`
	Secret string
}

type ServerConfig struct {
	Port int `default:"8080"`
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

	// articleService, err := article.NewFromDSN(config.Articles.DSN, config.Articles.Secret)
	// if err != nil {
	// 	log.Fatal("cannot create article service: ", err)
	// }

	store, err := postgres.NewStore(config.Users.DSN)
	if err != nil {
		log.Fatal("cannot create user store: ", err)
	}

	userServer, err := user.NewHTTP(store, []byte(config.Users.Secret))
	if err != nil {
		log.Fatal("cannot create user service: ", err)
	}

	// Setup API endpoints
	// router.Mount("/articles", articleService)
	router.Mount("/users", userServer)
	// router.Mount("/profiles", userService)

	log.Println("start listening on", server.Addr)
	log.Fatal(server.ListenAndServe())
}
