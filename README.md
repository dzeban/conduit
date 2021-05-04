# conduit

Real world app (Medium clone) backend implemented in Go

# Get started

    make up

This will create Postgres db, apply migrations and run API server - all via
docker-compose.

# TODO

[x] Make middleware for user auth to simplify handlers
    [x] User handlers
    [x] Article handlers
[x] Move service logic out of handlers to the service
[x] mock: Create test users in mock.NewUserStore constructor, expose test users
    in mock package
[x] Merge user/ subpackages into one "user" package. So we can avoid stupid
    types like service.Service and server.Server

[x] Implement profile service (follow/unfollow). It's required to implement
    articles feed.
[x] Migrate service/article to "articles" package. Add article methods to
    postgres store.
[x] Delete ./service/ and ./store/
[x] Fix cli util

[x] Create predefined errors. Merge all .../error.go into app/error.go

[ ] Write tests

    conduit has 3 types of tests:

    1. Unit tests to check business logic. In these unit tests we aim to cover
       as much as possible. It's tested against mock repository. Unit testing is
       performed:
        * In services (e.g. `user/login_test.go`, `article/list_test.go`).
        * In handlers (e.g. `profile/http_test.go`, `article/http_test.go`).
          Here it's checking that user gets response according to the API spec.
          In these tests we check correct HTTP status and error in the response.
          It is tested against mock repository.
        * In supporting code packages (e.g. `password/hash_test.go`, `jwt/token_test.go`).
    3. Integration tests in postgres package (`postgres/article_test.go`) to
       check that SQL queries are working with real Postgres. These tests are
       somewhat expensive, therefore they test the happy path.
    3. System tests (or functional, or end-to-end, you name it) that check the
       full user journey. They require full stack including database, they query
       API one or multiple times and check that the whole scenario succeeded.
       Scenario is some meaningful to user sequence of actions like "login is
       working", "article creation is working", "after following user I can see
       its articles in my feed", etc. These could've been implemented using BDD
       but in my opinion it's overkill so it's testing using simple `go test`.
       System tests live in its own package to avoid peeking into backend types
       for response checking. It's executed only when CONDUIT_TEST_API
       environment variable is set. I've discarded using build tags for this
       because with build tags most of the go tooling, e.g. linter, is not
       working.

[ ] Rewrite HTTP statuses to match API spec
[ ] Configure linters
[ ] Move to Gin
[ ] Remove down migrations
[ ] Move migrations to the store
[ ] Get rid of sql builder where it's not needed (except update methods)
[ ] Make custom lower snake case marshaller in local json package and remove
    json struct tags because they spoil types.
    https://gist.github.com/Rican7/39a3dc10c1499384ca91
[x] Remove vendoring

Think about refactoring inspired by WTFDial:

- Move service logic to the postgres package and merge it with queries. So
  service type is inside the postgres package.
- Instead of passing user into service method pass the context. And also use
  this context for postgres queries.
