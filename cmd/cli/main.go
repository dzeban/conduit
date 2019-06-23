package main

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"

	"github.com/dzeban/conduit/cmd/cli/state"
	"github.com/dzeban/conduit/cmd/cli/user"
)

func main() {
	// cli
	// > login username=test password=test
	// > register username=test password=test email=test@example.com
	// token is persisted in session
	// > articles
	// > articles/slug
	// >

	cli := ishell.NewWithConfig(&readline.Config{
		Prompt: fmt.Sprintf("(%s) > ", state.CurrentUsername),
	})

	// user commands
	{
		userCmd := &ishell.Cmd{
			Name: "user",
			Help: "Interact with user - login, register, etc.",
		}

		userCmd.AddCmd(&ishell.Cmd{
			Name:                "register",
			Help:                "Register new user",
			Func:                user.Register,
			CompleterWithPrefix: user.RegisterComplete,
		})

		userCmd.AddCmd(&ishell.Cmd{
			Name:                "login",
			Help:                "Login existing user",
			Func:                user.Login,
			CompleterWithPrefix: user.LoginComplete,
		})

		userCmd.AddCmd(&ishell.Cmd{
			Name: "get",
			Help: "Get current user",
			Func: user.Get,
		})

		cli.AddCmd(userCmd)
	}

	cli.SetHistoryPath(".conduit_history")
	cli.Run()
}
