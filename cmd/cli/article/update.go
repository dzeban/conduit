package article

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/cmd/cli/debug"
)

var UpdateOpts = []string{
	"slug",
	"title",
	"description",
	"body",
}

func Update(c *ishell.Context) {
	var slug string
	a := app.Article{}

	for _, opt := range c.Args {
		kv := strings.Split(opt, "=")
		if len(kv) != 2 {
			fmt.Println("invalid option format, need k=v")
			return
		}

		switch kv[0] {
		case "slug":
			slug = kv[1]
		case "title":
			a.Title = kv[1]
		case "description":
			a.Description = kv[1]
		case "body":
			a.Body = kv[1]
		}
	}

	if slug == "" {
		fmt.Println("slug is required")
		return
	}

	ar := app.ArticleUpdateRequest{
		Article: app.ArticleUpdateRequestData{
			Title:       a.Title,
			Description: a.Description,
			Body:        a.Body,
		},
	}

	url := fmt.Sprintf("http://localhost:8080/articles/%s", slug)
	_, _, err := debug.MakeAuthorizedRequestWithDump("PUT", url, ar)
	if err != nil {
		fmt.Println(err)
		return
	}
}
