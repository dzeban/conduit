.PHONY: up migrate

up:
	docker-compose up -d

down:
	docker-compose down

migrate:
	migrate -source file://migrations -database postgresql://postgres:postgres@localhost:5432/conduit?sslmode=disable up
