.PHONY: all build clean run

all: build

build:
	docker compose build

clean:
	docker compose down --rmi all

run:
	docker compose up
