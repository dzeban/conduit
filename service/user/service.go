package user

import (
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/db"
)

// Service provides a service for interacting with user accounts
type Service struct {
	db     *sqlx.DB
	router *mux.Router
	secret []byte
}

// NewService is a constructor for a Service
func NewService(DSN string, secret string) (*Service, error) {
	db, err := db.ConnectLoop("postgres", DSN, 1*time.Minute)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to users db")
	}

	router := mux.NewRouter().StrictSlash(true)

	s := &Service{db: db, router: router, secret: []byte(secret)}
	router.HandleFunc("/users/", s.HandleUserRegister).Methods("POST")
	router.HandleFunc("/users/login", s.HandleUserLogin).Methods("POST")
	router.HandleFunc("/users/", s.jwtAuthHandler(s.HandleUserGet)).Methods("GET")
	router.HandleFunc("/users/", s.jwtAuthHandler(s.HandleUserUpdate)).Methods("PUT")

	router.HandleFunc("/profiles/{username}", s.HandleProfileGet).Methods("GET")
	router.HandleFunc("/profiles/{username}/follow", s.jwtAuthHandler(s.HandleProfileFollow)).Methods("POST")
	router.HandleFunc("/profiles/{username}/unfollow", s.jwtAuthHandler(s.HandleProfileUnfollow)).Methods("POST")

	return s, nil
}
