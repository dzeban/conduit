package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestRegisterAndLogin(t *testing.T) {
	shouldSkip(t)

	cli := &http.Client{
		Timeout: time.Second * 1,
	}

	name := randString(10)
	email := name + "@example.com"
	password := name

	// Register
	apitest.New().
		EnableNetworking(cli).
		Post("http://localhost:8080/users/").
		JSON(fmt.Sprintf(`{"user": {"username": "%s", "email": "%s", "password": "%s"}}`, name, email, password)).
		Expect(t).
		Status(201).
		Assert(jsonpath.Present("user.token")).
		End()

	// Login
	apitest.New().
		EnableNetworking(cli).
		Post("http://localhost:8080/users/login").
		JSON(fmt.Sprintf(`{"user": {"email": "%s", "password": "%s"}}`, email, password)).
		Expect(t).
		Status(200).
		Assert(jsonpath.Present("user.token")).
		Assert(jsonpath.Equal("user.username", name)).
		Assert(jsonpath.Equal("user.email", email)).
		End()
}
