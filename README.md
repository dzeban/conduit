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
[ ] Migrate service/article to "articles" package. Add article methods to
    postgres store.
[ ] Delete ./service/ and ./store/
[ ] Fix cli util

[ ] Create predefined errors. Merge all .../error.go into app/error.go

[ ] Make custom lower snake case marshaller in local json package and remove
    json struct tags because they spoil types.
    https://gist.github.com/Rican7/39a3dc10c1499384ca91
[ ] Rewrite HTTP statuses to match API spec
[ ] Configure linters
[ ] Move to Gin
[ ] Remove down migrations
[ ] Move migrations to the store
[ ] Get rid of sql builder where it's not needed (except update methods)
[x] Remove vendoring

Think about refactoring inspired by WTFDial:

- Move service logic to the postgres package and merge it with queries. So
  service type is inside the postgres package.
- Instead of passing user into service method pass the context. And also use
  this context for postgres queries.
