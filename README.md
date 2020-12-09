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
[ ] Move service logic out of handlers to the service
[ ] Make custom lower snake case marshaller in local json package and remove
    json struct tags because they spoil types.
    https://gist.github.com/Rican7/39a3dc10c1499384ca91
[ ] Rewrite HTTP statuses to match API spec
[ ] Merge article service, store, handlers, requests into article package?
[ ] Create predefined errors
[ ] Move to Gin
[ ] Remove down migrations
[ ] Move migrations to the store
[ ] Remove vendoring
