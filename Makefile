## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	docker-compose run migrate create -seq -ext=.sql -dir=./migrations ${name}
