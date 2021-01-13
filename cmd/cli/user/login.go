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

var LoginOpts = []string{
	"email",
	"password",
}

func Login(c *ishell.Context) {
	// Construct login user request
	req := User.LoginRequest{}

	for _, opt := range c.Args {
		kv := strings.Split(opt, "=")
		if len(kv) != 2 {
			fmt.Println("invalid option format, need k=v")
			return
		}

		switch kv[0] {
		case "email":
			req.User.Email = kv[1]
		case "password":
			req.User.Password = kv[1]
		}
	}

	resp, body, err := debug.MakeRequestWithDump("POST", "http://localhost:8080/users/login", req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Update current state if request was successful
	if resp.StatusCode == http.StatusOK {
		var loginUser userRequest
		err = json.Unmarshal(body, &loginUser)
		if err != nil {
			fmt.Println(err)
			return
		}

		state.CurrentUsername = loginUser.User.Name
		state.CurrentToken = loginUser.User.Token
		fmt.Println("[update current user and token]")

		c.SetPrompt(fmt.Sprintf("(%s) > ", state.CurrentUsername))
	}
}
