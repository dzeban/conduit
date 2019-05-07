# conduit

Real world app (Medium clone) backend implemented in Go

# Get started


## Database setup

1. Wind up database with

    make up

This will start Postgres in container via docker-compose.

2. Create database `conduit` in Postgres

    CREATE DATABASE conduit

3. Apply migration

    make migrate

# TODO

- [x] Make error values, add String(), remove hardcoded string comparison in
  tests
- [x] Replace hardcoded integer status codes to net/http constants
- [x] Rewrite server configuration
- [ ] Add logging middleware
- [x] Fix error status codes
- [ ] Add "exp" claim to JWT, validate it
- [ ] Use Postgres implementation in server tests
- [ ] Reorganize packages - remove mock package, create `article`, `user`, etc.
  packages from postgres; move `app` package content to the dedicated packages,
  remove interfaces.
