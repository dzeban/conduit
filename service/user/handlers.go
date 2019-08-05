package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

var (
	ErrJWTNoAuthorizationHeader = errors.New("no Authorization header")
	ErrJWTNoSignedClaim         = errors.New("token does not have signed claim")
	ErrJWTNoSubClaim            = errors.New("no sub claim")
)

// ServeHTTP implements http.handler interface and uses router ServeHTTP method
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value("email")
	if val == nil {
		http.Error(w, ServerError(nil, "no email in context"), http.StatusUnauthorized)
		return
	}

	email, ok := val.(string)
	if !ok {
		http.Error(w, ServerError(nil, "invalid auth email"), http.StatusUnauthorized)
		return
	}

	if email == "" {
		http.Error(w, ServerError(nil, "empty auth email"), http.StatusUnauthorized)
		return
	}

	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req app.UserRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, ServerError(err, "failed to decode request"), http.StatusBadRequest)
		return
	}

	err = req.User.ValidateForUpdate()
	if err != nil {
		http.Error(w, ServerError(err, "failed to validate request"), http.StatusBadRequest)
		return
	}

	if req.User.Email != "" && req.User.Email != email {
		http.Error(w, ServerError(err, "not allowed to update other user"), http.StatusForbidden)
		return
	}

	user, err := s.Update(email, req.User)
	if err == app.ErrUserNotFound {
		http.Error(w, ServerError(err, "no such user"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, ServerError(err, "failed to update user"), http.StatusInternalServerError)
		return
	}

	// Generate new JWT because user was updated
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    user.Email,
		"signed": true,
	})

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		http.Error(w, ServerError(err, "failed to create token"), http.StatusInternalServerError)
		return
	}

	user.Token = tokenString

	// Prepare and send reply with user data, including token
	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		http.Error(w, ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}

func (s *Service) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req app.UserRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, ServerError(err, "failed to decode request"), http.StatusBadRequest)
		return
	}

	err = req.User.ValidateForRegister()
	if err != nil {
		http.Error(w, ServerError(err, "failed to validate request"), http.StatusBadRequest)
		return
	}

	user, err := s.Register(req.User)
	if err == app.ErrUserExists {
		http.Error(w, ServerError(err, "failed to register user"), http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, ServerError(err, "failed to register user"), http.StatusInternalServerError)
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    user.Email,
		"signed": true,
	})

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		http.Error(w, ServerError(err, "failed to create token"), http.StatusInternalServerError)
		return
	}

	user.Token = tokenString

	// Prepare and send reply with user data, including token
	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		http.Error(w, ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonUser)
}

func (s *Service) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req app.UserRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, ServerError(err, "failed to decode request"), http.StatusBadRequest)
		return
	}

	err = req.User.ValidateForLogin()
	if err != nil {
		http.Error(w, ServerError(err, "failed to validate request"), http.StatusBadRequest)
		return
	}

	user, err := s.Login(req.User)
	if err != nil {
		http.Error(w, ServerError(err, "failed to login"), http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    user.Email,
		"signed": true,
	})

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		http.Error(w, ServerError(err, "failed to create token"), http.StatusInternalServerError)
		return
	}

	user.Token = tokenString

	// Prepare and send reply with user data, including token
	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		http.Error(w, ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.Write(jsonUser)
}

// HandleUserGet gets the currently logged-in user. Requires authentication.
func (s *Service) HandleUserGet(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value("email")
	if val == nil {
		http.Error(w, ServerError(nil, "no email in context"), http.StatusUnauthorized)
		return
	}

	email, ok := val.(string)
	if !ok {
		http.Error(w, ServerError(nil, "invalid auth email"), http.StatusUnauthorized)
		return
	}

	if email == "" {
		http.Error(w, ServerError(nil, "empty auth email"), http.StatusUnauthorized)
		return
	}

	user, err := s.Get(email)
	if err == app.ErrUserNotFound {
		http.Error(w, ServerError(err, "failed to get user"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, ServerError(err, "failed to get user"), http.StatusInternalServerError)
		return
	}

	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		http.Error(w, ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.Write(jsonUser)
}

func (s *Service) HandleProfileGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	profile, err := s.store.Profile(username)
	if err == app.ErrUserNotFound {
		http.Error(w, ServerError(err, "no such user"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, ServerError(err, "failed to get profile"), http.StatusInternalServerError)
		return
	}

	jsonProfile, err := json.Marshal(app.ProfileResponse{Profile: *profile})
	if err != nil {
		http.Error(w, ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.Write(jsonProfile)
}

func (s *Service) HandleProfileFollow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, ServerError(nil, "no email in context"), http.StatusUnauthorized)
		return
	}

	// Query current user to get name
	follower, err := s.Get(email)
	if err == app.ErrUserNotFound {
		http.Error(w, ServerError(err, "no such user"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, ServerError(err, "failed to get user"), http.StatusInternalServerError)
		return
	}

	// Query profile to follow to ensure it exists
	profile, err := s.store.Profile(username)
	if err == app.ErrUserNotFound {
		http.Error(w, ServerError(err, "no such profile"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, ServerError(err, "failed to get profile"), http.StatusInternalServerError)
		return
	}

	err = s.store.Follow(follower.Name, username)
	if err != nil {
		http.Error(w, ServerError(err, "failed to follow user"), http.StatusInternalServerError)
		return
	}

	// Prepare and send reply with profile data
	jsonUser, err := json.Marshal(app.ProfileResponse{Profile: *profile})
	if err != nil {
		http.Error(w, ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}

func (s *Service) HandleProfileUnfollow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, ServerError(nil, "no email in context"), http.StatusUnauthorized)
		return
	}

	// Query current user to get name
	follower, err := s.Get(email)
	if err == app.ErrUserNotFound {
		http.Error(w, ServerError(err, "no such user"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, ServerError(err, "failed to get user"), http.StatusInternalServerError)
		return
	}

	// Query profile to follow to ensure it exists
	profile, err := s.store.Profile(username)
	if err == app.ErrUserNotFound {
		http.Error(w, ServerError(err, "no such profile"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, ServerError(err, "failed to get profile"), http.StatusInternalServerError)
		return
	}

	err = s.store.Unfollow(follower.Name, username)
	if err != nil {
		http.Error(w, ServerError(err, "failed to follow user"), http.StatusInternalServerError)
		return
	}

	// Prepare and send reply with profile data
	jsonUser, err := json.Marshal(app.ProfileResponse{Profile: *profile})
	if err != nil {
		http.Error(w, ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}

// jwtAuthHandler is a middleware that wraps next handler func with JWT token
// parsing and validation. It also stores authenticated user email into the
// context.
func (s *Service) jwtAuthHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader, ok := r.Header["Authorization"]
		if !ok {
			http.Error(w, ServerError(ErrJWTNoAuthorizationHeader, ""), http.StatusUnauthorized)
			return
		}

		claims, err := parseJWTClaimsFromHeader(authHeader[0], s.secret)
		if err != nil {
			http.Error(w, ServerError(err, "failed to parse jwt"), http.StatusBadRequest)
			return
		}

		if claims["signed"] != true {
			http.Error(w, ServerError(ErrJWTNoSignedClaim, ""), http.StatusUnauthorized)
			return
		}

		var sub interface{}
		if sub, ok = claims["sub"]; !ok {
			http.Error(w, ServerError(ErrJWTNoSubClaim, ""), http.StatusUnauthorized)
			return
		}

		// Store auth subject (email) to the context
		authCtx := context.WithValue(r.Context(), "email", sub)

		next(w, r.WithContext(authCtx))
	}
}

func (s *Service) loggerHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}
