.PHONY: all build clean run test

all: build

build:
	docker compose build

clean:
	docker compose down --rmi all

run:
	docker compose up

test:
	docker compose run --rm app go test ./...
