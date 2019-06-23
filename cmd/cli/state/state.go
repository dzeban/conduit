package state

import "net/http"

var (
	CurrentUsername string
	CurrentToken    string
)

var Client = &http.Client{}
