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

