package main

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"

	"github.com/dzeban/conduit/cmd/cli/profile"
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

		userCmd.AddCmd(&ishell.Cmd{
			Name:                "update",
			Help:                "Update current user",
			Func:                user.Update,
			CompleterWithPrefix: OptsCompleter(user.UpdateOpts),
		})

		cli.AddCmd(userCmd)
	}

	// profile commands
	{
		profileCmd := &ishell.Cmd{
			Name: "profile",
			Help: "Interact with profiles - get, follow, unfollow",
		}

		profileCmd.AddCmd(&ishell.Cmd{
			Name:                "get",
			Help:                "Get profile by name",
			Func:                profile.Get,
			CompleterWithPrefix: OptsCompleter(profile.GetOpts),
		})

		cli.AddCmd(profileCmd)
	}

	cli.SetHistoryPath(".conduit_history")
	cli.Run()
}
