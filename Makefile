.PHONY: build run test clean docker-up docker-down

build:
	go build -o bin/app ./cmd/app

run:
	go run ./cmd/app

test:
	go test -v ./...

clean:
	rm -rf bin/

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v