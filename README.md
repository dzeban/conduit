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

- [ ] Improve coverage for HandleArticles
- [ ] Make error values, add String(), remove hardcoded string comparison in
  tests
- [ ] Replace hardcoded integer status codes to net/http constants
- [ ] Rewrite server configuration with confita
- [ ] Add logging middleware
- [ ] Fix error status codes
- [ ] Add "exp" claim to JWT, validate it
