debug:
	docker-compose -f ./docker-compose.yaml -f ./docker-compose.debug.yaml up -d
down:
	docker-compose down
prod:
	docker-compose -f ./docker-compose.yaml up --build
