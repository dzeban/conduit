package profile

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"

	"github.com/dzeban/conduit/cmd/cli/debug"
)

var GetOpts = []string{
	"username",
}

func Get(c *ishell.Context) {
	if len(c.Args) < 1 {
		fmt.Println("username is required")
		return
	}

	kv := strings.Split(c.Args[0], "=")
	if len(kv) != 2 {
		fmt.Println("invalid option format, need k=v")
		return
	}

	username := kv[1]

	url := fmt.Sprintf("http://localhost:8080/profiles/%s", username)
	_, _, err := debug.MakeRequestWithDump("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
