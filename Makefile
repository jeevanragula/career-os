APP_NAME := careeros
COMPOSE := docker compose

.PHONY: setup up dev down stop restart logs logs-app logs-db status ps build test migrate db-shell reset clean fmt

# Recommended one-command local setup/start.
setup:
	@$(MAKE) up
	@echo ""
	@echo "CareerOS: http://localhost:8080"
	@echo "Logs:     make logs"

# Build and start in background. Migrations run automatically at app startup.
up:
	$(COMPOSE) up --build -d

# Build and start in foreground.
dev:
	$(COMPOSE) up --build

# Stop containers but preserve database data.
down:
	$(COMPOSE) down

stop: down

restart:
	$(COMPOSE) down
	$(COMPOSE) up --build -d

logs:
	$(COMPOSE) logs -f --tail=200 careeros

logs-app: logs

logs-db:
	$(COMPOSE) logs -f --tail=200 postgres

status:
	$(COMPOSE) ps

ps: status

build:
	$(COMPOSE) build

test:
	go test ./...

migrate:
	$(COMPOSE) up -d postgres
	$(COMPOSE) run --rm careeros /app/career-migrate

db-shell:
	$(COMPOSE) exec postgres psql -U careeros -d careeros

# WARNING: deletes the local Postgres volume and ALL local CareerOS data.
reset:
	$(COMPOSE) down -v --remove-orphans
	$(COMPOSE) up --build -d

clean:
	$(COMPOSE) down -v --remove-orphans

fmt:
	gofmt -w cmd internal
