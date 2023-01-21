.PHONY: confirm
confirm:
	@/bin/echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	docker-compose run migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	docker-compose run migrate up