build:
	docker compose build

up:
	docker compose up

down:
	docker compose stop

cleanup:
	docker compose down -v