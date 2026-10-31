include .env
export

export PROJECT_ROOT = ${shell pwd}

env-up:
	@docker compose up -d auth-postgres
# 	ToDo

env-down:
	@docker compose down auth-postgres
# 	ToDo

env-cleanup:
	@read -p "Are you sure you want to remove the database volume? This action is irreversible. [y/n]: " confirm; \
	if [ "$$confirm" = "y" ]; then \
		docker compose down auth-postgres -v && \
		echo "Env cleaned up."; \
	else \
		echo "Cleanup aborted."; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Error: Please provide a migration name using 'make migrate-create seq=<migration_seq>'"; \
		exit 1; \
	fi; \
	docker compose run --rm --user $$(id -u):$$(id -g) auth-service-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Error: Please provide an action (up or down) using 'make migrate-action action=<up|down>'"; \
		exit 1; \
	fi; \
	docker compose run --rm auth-service-migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		"$(action)"

auth-service-run:
	@go run ./cmd/auth-service/main.go