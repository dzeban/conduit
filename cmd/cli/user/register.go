package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/abiosoft/ishell"

	"github.com/dzeban/conduit/cmd/cli/debug"
	"github.com/dzeban/conduit/cmd/cli/state"
)

var RegisterOpts = []string{
	"username",
	"email",
	"bio",
	"image",
	"password",
}

func Register(c *ishell.Context) {
	// Construct register user request
	u := user{}

	for _, opt := range c.Args {
		kv := strings.Split(opt, "=")
		if len(kv) != 2 {
			fmt.Println("invalid option format, need k=v")
			return
		}

		switch kv[0] {
		case "username":
			u.Name = kv[1]
		case "email":
			u.Email = kv[1]
		case "bio":
			u.Bio = kv[1]
		case "image":
			u.Image = kv[1]
		case "password":
			u.Password = kv[1]
		}
	}

	ur := userRequest{
		User: u,
	}

	resp, body, err := debug.MakeRequestWithDump("POST", "http://localhost:8080/users/", ur)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Update current state if request was successful
	if resp.StatusCode == http.StatusCreated {
		var registeredUser userRequest
		err = json.Unmarshal(body, &registeredUser)
		if err != nil {
			fmt.Println(err)
			return
		}

		state.CurrentUsername = registeredUser.User.Name
		state.CurrentToken = registeredUser.User.Token
		fmt.Println("[update current user and token]")

		c.SetPrompt(fmt.Sprintf("(%s) > ", state.CurrentUsername))
	}
}
