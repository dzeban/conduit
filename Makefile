.PHONY: up down stop migrate

run: test build up
	./conduit

test:
	go test ./...

build:
	go build ./cmd/conduit

install:
	go install ./cmd/conduit

up:
	docker-compose up --build

down:
	docker-compose down

start:
	docker-compose up

stop:
	docker-compose stop

swagger-up:
	 docker run -p 8888:8080 --name swagger-conduit -e SWAGGER_JSON=/api/swagger.json -v $$(pwd)/api:/api -d swaggerapi/swagger-ui
	 xdg-open http://localhost:8888

swagger-down:
	docker stop swagger-conduit
	docker rm swagger-conduit
