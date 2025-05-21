.PHONY: build run test docker-build docker-run clean

BINARY_NAME=traffic-controller
MAIN_FILE=cmd/server/main.go

build:
	go build -o $(BINARY_NAME) $(MAIN_FILE)

run:
	go run $(MAIN_FILE)

test:
	go test ./... -v

docker-build:
	docker-compose build

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

clean:
	rm -f $(BINARY_NAME)
	go clean

init:
	cp .env.example .env
	mkdir -p scripts
	mkdir -p scripts/init.sql