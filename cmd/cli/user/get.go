package user

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/dzeban/conduit/cmd/cli/debug"
)

func Get(c *ishell.Context) {
	_, _, err := debug.MakeAuthorizedRequestWithDump("GET", "http://localhost:8080/users/", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
