package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/dzeban/conduit/cmd/cli/debug"
	"github.com/dzeban/conduit/cmd/cli/state"
	User "github.com/dzeban/conduit/user"
)

var UpdateOpts = []string{
	"username",
	"email",
	"bio",
	"image",
	"password",
}

func Update(c *ishell.Context) {
	// Construct register user request
	req := User.UpdateRequest{}

	for _, opt := range c.Args {
		kv := strings.Split(opt, "=")
		if len(kv) != 2 {
			fmt.Println("invalid option format, need k=v")
			return
		}

		switch kv[0] {
		case "username":
			req.User.Name = kv[1]
		case "email":
			req.User.Email = kv[1]
		case "bio":
			req.User.Bio = kv[1]
		case "image":
			req.User.Image = kv[1]
		case "password":
			req.User.Password = kv[1]
		}
	}

	resp, body, err := debug.MakeAuthorizedRequestWithDump("PUT", "http://localhost:8080/users/", req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Update current state if request was successful
	if resp.StatusCode == http.StatusOK {
		var updatedUser userRequest
		err = json.Unmarshal(body, &updatedUser)
		if err != nil {
			fmt.Println(err)
			return
		}

		state.CurrentUsername = updatedUser.User.Name
		state.CurrentToken = updatedUser.User.Token
		fmt.Println("[update current user and token]")

		c.SetPrompt(fmt.Sprintf("(%s) > ", state.CurrentUsername))
	}
}
