build:
	docker compose build

up:
	docker compose up

down:
	docker compose stop

cleanup:
	docker compose down -v

getTemperature:
	curl 'localhost:3100/get?key=temperature'

getHumidity:
	curl 'localhost:3100/get?key=humidity'

setTemperature:
	curl 'localhost:3100/set?key=temperature&value=15'

setHumidity:
	curl 'localhost:3100/set?key=humidity&value=50'

