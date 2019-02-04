.PHONY: up down stop migrate

up:
	docker-compose up -d

down:
	docker-compose down

stop:
	docker-compose stop

migrate:
	migrations/migrate -source file://migrations -database postgresql://postgres:postgres@localhost:5432/conduit?sslmode=disable up

build: *.go
	go build

run: build
	./conduit
