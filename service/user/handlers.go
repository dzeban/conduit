package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/jwt"
)

// ServeHTTP implements http.handler interface and uses router ServeHTTP method
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		http.Error(w, app.ServerError(nil, "no user in context"), http.StatusUnauthorized)
		return
	}

	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req app.UserRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to decode request"), http.StatusBadRequest)
		return
	}

	err = req.User.ValidateForUpdate()
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to validate request"), http.StatusBadRequest)
		return
	}

	user, err := s.Update(currentUser.Email, req.User)
	if err == app.ErrUserNotFound {
		http.Error(w, app.ServerError(err, "no such user"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, app.ServerError(err, "failed to update user"), http.StatusInternalServerError)
		return
	}

	// Generate new JWT because user was updated
	token, err := jwt.New(user, s.secret)
	if err != nil {
		http.Error(w, app.ServerError(err, ""), http.StatusInternalServerError)
		return
	}

	user.Token = token

	// Prepare and send reply with user data, including token
	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
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
		http.Error(w, app.ServerError(err, "failed to decode request"), http.StatusBadRequest)
		return
	}

	err = req.User.ValidateForRegister()
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to validate request"), http.StatusBadRequest)
		return
	}

	user, err := s.Register(req.User)
	if err == app.ErrUserExists {
		http.Error(w, app.ServerError(err, "failed to register user"), http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, app.ServerError(err, "failed to register user"), http.StatusInternalServerError)
		return
	}

	// Generate JWT
	token, err := jwt.New(user, s.secret)
	if err != nil {
		http.Error(w, app.ServerError(err, ""), http.StatusInternalServerError)
		return
	}

	user.Token = token

	// Prepare and send reply with user data, including token
	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
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
		http.Error(w, app.ServerError(err, "failed to decode request"), http.StatusBadRequest)
		return
	}

	err = req.User.ValidateForLogin()
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to validate request"), http.StatusBadRequest)
		return
	}

	user, err := s.Login(req.User)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to login"), http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := jwt.New(user, s.secret)
	if err != nil {
		http.Error(w, app.ServerError(err, ""), http.StatusInternalServerError)
		return
	}

	user.Token = token

	// Prepare and send reply with user data, including token
	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.Write(jsonUser)
}

// HandleUserGet gets the currently logged-in user. Requires authentication.
func (s *Service) HandleUserGet(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		http.Error(w, app.ServerError(nil, "no user in context"), http.StatusUnauthorized)
		return
	}

	user, err := s.Get(currentUser.Email)
	if err == app.ErrUserNotFound {
		http.Error(w, app.ServerError(err, "failed to get user"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, app.ServerError(err, "failed to get user"), http.StatusInternalServerError)
		return
	}

	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.Write(jsonUser)
}

func (s *Service) HandleProfileGet(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	profile, err := s.store.Profile(username)
	if err == app.ErrUserNotFound {
		http.Error(w, app.ServerError(err, "no such user"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, app.ServerError(err, "failed to get profile"), http.StatusInternalServerError)
		return
	}

	jsonProfile, err := json.Marshal(app.ProfileResponse{Profile: *profile})
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.Write(jsonProfile)
}

func (s *Service) HandleProfileFollow(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		http.Error(w, app.ServerError(nil, "no user in context"), http.StatusUnauthorized)
		return
	}

	username := chi.URLParam(r, "username")

	// Query current user to get name
	follower, err := s.Get(currentUser.Email)
	if err == app.ErrUserNotFound {
		http.Error(w, app.ServerError(err, "no such user"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, app.ServerError(err, "failed to get user"), http.StatusInternalServerError)
		return
	}

	// Query profile to follow to ensure it exists
	profile, err := s.store.Profile(username)
	if err == app.ErrUserNotFound {
		http.Error(w, app.ServerError(err, "no such profile"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, app.ServerError(err, "failed to get profile"), http.StatusInternalServerError)
		return
	}

	err = s.store.Follow(follower.Name, username)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to follow user"), http.StatusInternalServerError)
		return
	}

	// Prepare and send reply with profile data
	jsonUser, err := json.Marshal(app.ProfileResponse{Profile: *profile})
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}

func (s *Service) HandleProfileUnfollow(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		http.Error(w, app.ServerError(nil, "no user in context"), http.StatusUnauthorized)
		return
	}

	// Query current user to get name
	follower, err := s.Get(currentUser.Email)
	if err == app.ErrUserNotFound {
		http.Error(w, app.ServerError(err, "no such user"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, app.ServerError(err, "failed to get user"), http.StatusInternalServerError)
		return
	}

	// Query profile to follow to ensure it exists
	profile, err := s.store.Profile(username)
	if err == app.ErrUserNotFound {
		http.Error(w, app.ServerError(err, "no such profile"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, app.ServerError(err, "failed to get profile"), http.StatusInternalServerError)
		return
	}

	err = s.store.Unfollow(follower.Name, username)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to follow user"), http.StatusInternalServerError)
		return
	}

	// Prepare and send reply with profile data
	jsonUser, err := json.Marshal(app.ProfileResponse{Profile: *profile})
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}
