up:
	docker compose up --build

up-detached:
	docker compose up --build -d

down:
	docker compose down

reset:
	docker compose down -v
	docker compose up --build

logs:
	docker compose logs -f careeros

test:
	go test ./...

fmt:
	gofmt -w cmd internal
