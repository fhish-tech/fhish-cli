BINARY_NAME=fhish

build:
	go build -o bin/$(BINARY_NAME) main.go

install:
	go install

clean:
	rm -rf bin/
	go clean

COMPOSE_FILE = docker/docker-compose.yml

docker-up: build
	docker compose -f $(COMPOSE_FILE) up --build -d

docker-down:
	docker compose -f $(COMPOSE_FILE) down -v

docker-logs:
	docker compose -f $(COMPOSE_FILE) logs -f

docker-status:
	docker compose -f $(COMPOSE_FILE) ps

docker-verify:
	docker compose -f $(COMPOSE_FILE) --profile verify run --rm verifier

test:
	go test ./...

.PHONY: build install clean run test docker-up docker-down docker-logs docker-status docker-verify
