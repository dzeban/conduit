package transport

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/dzeban/conduit/app"
)

type ErrorResponse struct {
	Errors Errors `json:"errors"`
}

type Errors struct {
	Body []string `json:"body"`
}

type HandlerWithError func(http.ResponseWriter, *http.Request) error

func WithError(h HandlerWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Invoke handler and get its error
		err := h(w, r)
		if err != nil {
			// If we got error, unwrap it to the app.Error to properly serialize
			var e app.Error
			ok := errors.As(err, &e)
			if !ok {
				log.Printf("invalid error %T from handler: %+v", err, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			switch e.Type {
			// Internal server errors are not returned to user, they are logged
			case app.ErrorTypeInternal:
				log.Printf("internal server error: %+v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return

			case app.ErrorTypeAuth:
				w.WriteHeader(http.StatusUnauthorized)

			default:
				w.WriteHeader(http.StatusUnprocessableEntity)
			}

			w.Header().Set("Content-Type", "application/json")

			err = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: Errors{
					Body: []string{err.Error()},
				},
			})
			if err != nil {
				log.Printf("failed to marshal error response: %+v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	}
}
