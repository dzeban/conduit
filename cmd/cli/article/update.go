package article

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"

	Article "github.com/dzeban/conduit/article"
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
	req := Article.UpdateRequest{}

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
			req.Article.Title = kv[1]
		case "description":
			req.Article.Description = kv[1]
		case "body":
			req.Article.Body = kv[1]
		}
	}

	if slug == "" {
		fmt.Println("slug is required")
		return
	}

	url := fmt.Sprintf("http://localhost:8080/articles/%s", slug)
	_, _, err := debug.MakeAuthorizedRequestWithDump("PUT", url, req)
	if err != nil {
		fmt.Println(err)
		return
	}
}
