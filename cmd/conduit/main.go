package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Server.Port),
		Handler: router,
	}

	log.Printf("using config: %#v\n", config)

	articleService, err := article.NewService(config.Articles.DSN)
	if err != nil {
		log.Fatal("cannot create articles service: ", err)
	}

	userService, err := user.NewService(config.Users.DSN, config.Users.Secret)
	if err != nil {
		log.Fatal("cannot create users service: ", err)
	}

	// Setup API endpoints
	router.PathPrefix("/articles").Handler(articleService)
	router.PathPrefix("/users").Handler(userService)

	log.Println("start listening on", server.Addr)
	log.Fatal(server.ListenAndServe())
}
