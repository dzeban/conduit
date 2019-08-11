package article

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/dzeban/conduit/cmd/cli/debug"
)

func List(c *ishell.Context) {
	_, _, err := debug.MakeRequestWithDump("GET", "http://localhost:8080/articles/", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
