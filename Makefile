include .env
export

env-up:
	@docker compose up -d tasks-mysql tasks-redis tasks-app

env-down:
	@docker compose down tasks-mysql tasks-redis tasks-app

env-cleanup:
	@docker compose down -v

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@docker compose run --rm --no-deps tasks-mysql-migrate create -ext sql -dir /migrations -seq $(name)

migrate-up:
	@docker compose run --rm tasks-mysql-migrate -path /migrations -database ${DATABASE_URL} up

migrate-down:
	@docker compose run --rm tasks-mysql-migrate -path /migrations -database ${DATABASE_URL} down 1

run:
	@set LOGGER_FOLDER=out\logs&&go mod tidy &&go run cmd/tasks/main.go