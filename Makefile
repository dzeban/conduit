.PHONY: run test build install up down start stop swagger-up swagger-down cli psql

run: test build up
	./conduit

test:
	go test ./...

build:
	go build ./cmd/conduit

install:
	go install ./cmd/conduit

restart:
	docker-compose down
	docker-compose up --build

up:
	docker-compose up --build

down:
	docker-compose down

start:
	docker-compose up --build -d

stop:
	docker-compose stop

swagger-up:
	 docker run -p 8888:8080 --name swagger-conduit -e SWAGGER_JSON=/api/swagger.json -v $$(pwd)/api:/api -d swaggerapi/swagger-ui
	 xdg-open http://localhost:8888

swagger-down:
	docker stop swagger-conduit
	docker rm swagger-conduit

integration-test:
	docker-compose -f docker-compose.test.yml down
	docker-compose -f docker-compose.test.yml up --build

integration-test-down:
	docker-compose -f docker-compose.test.yml down

integration-test-psql:
	docker-compose -f docker-compose.test.yml exec postgres psql -U test

cli:
	go build ./cmd/cli

psql:
	docker-compose exec postgres psql -U conduit
