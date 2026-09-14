include .env
export

export PROJECT_ROOT=$(shell pwd)

ps:
	@docker compose ps

env-up:
	@docker compose up -d honey_postgres

env-down:
	@docker compose down honey_postgres

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует параметр seq"; \
		exit 1; \
	fi; \
	docker compose run --rm honey_postgres_migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует параметр action"; \
		exit 1; \
	fi; \
	docker compose run --rm honey_postgres_migrate \
		-path /migrations \
		-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@honey_postgres:5432/$(POSTGRES_DB)?sslmode=disable \
		"$(action)"

migrate-up:
	@$(MAKE) --no-print-directory migrate-action action=up

migrate-down:
	@$(MAKE) --no-print-directory migrate-action action=down

env-port-forward:
	@docker compose up -d port_forwarder

env-port-close:
	@docker compose down port_forwarder