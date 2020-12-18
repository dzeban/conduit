package transport

import (
	"encoding/json"
	"errors"
)

var (
	ErrorUnmarshal = errors.New("failed to unmarshal request")
	ErrorMarshal   = errors.New("failed to marshal response")
	ErrorJWT       = errors.New("JWT error")
)

type ErrorResponse struct {
	Errors Errors `json:"errors"`
}

type Errors struct {
	Body []string `json:"body"`
}

func ServerError(errors ...error) string {
	errorStrings := make([]string, 0, len(errors))

	for _, e := range errors {
		errorStrings = append(errorStrings, e.Error())
	}

	resp, err := json.Marshal(ErrorResponse{
		Errors: Errors{
			Body: errorStrings,
		},
	})
	if err != nil {
		// TODO: log error
		panic(err)
	}

	return string(resp)
}
