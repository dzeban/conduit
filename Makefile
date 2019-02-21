.PHONY: up down stop migrate

run: build up
	./conduit

build: *.go
	go build

up:
	docker-compose up -d

down:
	docker-compose down

stop:
	docker-compose stop

migrate:
	migrations/migrate -source file://migrations -database postgresql://postgres:postgres@localhost:5432/conduit?sslmode=disable up

