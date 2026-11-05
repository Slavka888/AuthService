export PROJECT_ROOT = ${shell pwd}

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
	docker compose run --rm --user $$(id -u):$$(id -g) auth-service-migrate \
		-path /migrations \
		-database "sqlite:///storage/auth.db" \
		"$(action)"

auth-service-run:
	@go run ./cmd/auth-service/main.go