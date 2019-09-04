package article

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/dzeban/conduit/cmd/cli/debug"
)

var FeedOpts = []string{
	"limit",
	"offset",
}

func Feed(c *ishell.Context) {
	var params [][]string
	for _, opt := range c.Args {
		kv := strings.Split(opt, "=")
		if len(kv) != 2 {
			fmt.Println("invalid option format, need k=v")
			return
		}

		params = append(params, kv)
	}

	_, _, err := debug.MakeAuthorizedRequestWithDump("GET", "http://localhost:8080/articles/feed", nil, params...)
	if err != nil {
		fmt.Println(err)
		return
	}
}
