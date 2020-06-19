.PHONY: help up start up-pg up-redis status logs restart clean pg

up:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

start:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --detach

up-pg:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml run --publish=127.0.0.1:5432:5432 postgres

up-redis:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml run --publish=127.0.0.1:6379:6379 redis

status:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml ps

logs:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs --tail=100

restart:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml stop
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --detach

clean:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down

pg:
	docker exec -it hotdeals_postgres-dev sh

