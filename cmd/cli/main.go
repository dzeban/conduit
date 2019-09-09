package main

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"

	"github.com/dzeban/conduit/cmd/cli/article"
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

		profileCmd.AddCmd(&ishell.Cmd{
			Name:                "follow",
			Help:                "Follow user by name",
			Func:                profile.Follow,
			CompleterWithPrefix: OptsCompleter(profile.FollowOpts),
		})

		profileCmd.AddCmd(&ishell.Cmd{
			Name:                "unfollow",
			Help:                "Unfollow user by name",
			Func:                profile.Unfollow,
			CompleterWithPrefix: OptsCompleter(profile.UnfollowOpts),
		})

		cli.AddCmd(profileCmd)
	}

	// article commands
	{
		articleCmd := &ishell.Cmd{
			Name: "article",
			Help: "Interact with articles - get, list, etc",
		}

		articleCmd.AddCmd(&ishell.Cmd{
			Name:                "get",
			Help:                "Get article by slug",
			Func:                article.Get,
			CompleterWithPrefix: OptsCompleter(article.GetOpts),
		})

		articleCmd.AddCmd(&ishell.Cmd{
			Name:                "list",
			Help:                "List articles",
			Func:                article.List,
			CompleterWithPrefix: OptsCompleter(article.ListOpts),
		})

		articleCmd.AddCmd(&ishell.Cmd{
			Name:                "feed",
			Help:                "Feed of articles",
			Func:                article.Feed,
			CompleterWithPrefix: OptsCompleter(article.FeedOpts),
		})

		articleCmd.AddCmd(&ishell.Cmd{
			Name:                "create",
			Help:                "Create article",
			Func:                article.Create,
			CompleterWithPrefix: OptsCompleter(article.CreateOpts),
		})

		cli.AddCmd(articleCmd)
	}

	cli.SetHistoryPath(".conduit_history")
	cli.Run()
}
