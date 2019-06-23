package main

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"

	"github.com/dzeban/conduit/cmd/cli/state"
	"github.com/dzeban/conduit/cmd/cli/user"
)

func main() {
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
			CompleterWithPrefix: OptsCompleter(user.RegisterOpts),
		})

		userCmd.AddCmd(&ishell.Cmd{
			Name:                "login",
			Help:                "Login existing user",
			Func:                user.Login,
			CompleterWithPrefix: OptsCompleter(user.LoginOpts),
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
