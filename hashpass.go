package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dzeban/conduit/password"
)

func generate(input string) {
	s, err := password.HashAndEncode(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s)
}

func compare(input, encodedHash string) {
	equals, err := password.Check(input, encodedHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(equals)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: %s <string> [<encoded hash>]")
		os.Exit(1)
	}

	if len(os.Args) == 2 {
		generate(os.Args[1])
	}

	if len(os.Args) == 3 {
		compare(os.Args[1], os.Args[2])
	}
}
