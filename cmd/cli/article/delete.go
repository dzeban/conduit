package article

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/dzeban/conduit/cmd/cli/debug"
)

var DeleteOpts = []string{
	"slug",
}

func Delete(c *ishell.Context) {
	if len(c.Args) < 1 {
		fmt.Println("slug is required")
		return
	}

	kv := strings.Split(c.Args[0], "=")
	if len(kv) != 2 {
		fmt.Println("invalid option format, need k=v")
		return
	}

	slug := kv[1]

	url := fmt.Sprintf("http://localhost:8080/articles/%s", slug)
	_, _, err := debug.MakeAuthorizedRequestWithDump("DELETE", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
