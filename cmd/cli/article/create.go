package article

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"

	Article "github.com/dzeban/conduit/article"
	"github.com/dzeban/conduit/cmd/cli/debug"
)

var CreateOpts = []string{
	"title",
	"description",
	"body",
}

func Create(c *ishell.Context) {
	req := Article.CreateRequest{}

	for _, opt := range c.Args {
		kv := strings.Split(opt, "=")
		if len(kv) != 2 {
			fmt.Println("invalid option format, need k=v")
			return
		}

		switch kv[0] {
		case "title":
			req.Article.Title = kv[1]
		case "description":
			req.Article.Description = kv[1]
		case "body":
			req.Article.Body = kv[1]
		}
	}

	_, _, err := debug.MakeAuthorizedRequestWithDump("POST", "http://localhost:8080/articles/", req)
	if err != nil {
		fmt.Println(err)
		return
	}
}
