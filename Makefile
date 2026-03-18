.PHONY: up down build logs restart clean

up:
	docker compose up --build -d

down:
	docker compose down

build:
	docker compose build

restart:
	docker compose restart app

clean:
	docker compose down -v
