.PHONY: up migrate

up:
	docker-compose up -d

down:
	docker-compose down

stop:
	docker-compose stop

migrate:
	migrate -source file://migrations -database postgresql://postgres:postgres@localhost:5432/conduit?sslmode=disable up
