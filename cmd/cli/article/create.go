package article

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/cmd/cli/debug"
)

var CreateOpts = []string{
	"title",
	"description",
	"body",
}

func Create(c *ishell.Context) {
	a := app.Article{}

	for _, opt := range c.Args {
		kv := strings.Split(opt, "=")
		if len(kv) != 2 {
			fmt.Println("invalid option format, need k=v")
			return
		}

		switch kv[0] {
		case "title":
			a.Title = kv[1]
		case "description":
			a.Description = kv[1]
		case "body":
			a.Body = kv[1]
		}
	}

	ar := app.ArticleCreateRequest{
		Article: a,
	}

	_, _, err := debug.MakeAuthorizedRequestWithDump("POST", "http://localhost:8080/articles/", ar)
	if err != nil {
		fmt.Println(err)
		return
	}
}
