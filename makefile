tailwindCmd = npx @tailwindcss/cli
tailwindMinify = --minify
prodFlags = -ldflags "-X internal/env.Environment=production"
out = ./bin/gorssagg

.PHONY: dev
dev:
	air

.PHONY: build-app
build-app:
	make tailwind
	make templ-generate
	go build ${prodFlags} -o ${out} ./cmd/main.go

.PHONY: build-app-dev
build-app-dev:
	make build-app prodFlags="" out=./tmp/main tailwindMinify=""
	templ generate --notify-proxy

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: tailwind
tailwind:
	$(tailwindCmd) -i ./main.css -o ./public/css/styles.css $(tailwindWatch) $(tailwindMinify)

.PHONY: tailwind-watch
tailwind-watch:
	make tailwind tailwindWatch=--watch

.PHONY: fetch-tailwind
fetch-tailwind:
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-arm64-musl

.PHONY: migrate-dev-up
migrate-dev-up:
	goose -dir ./sql/schema postgres "postgres://postgres:postgres@localhost:5432/gorssagg?sslmode=disable" up

.PHONY: migrate-dev-down
migrate-dev-down:
	goose -dir ./sql/schema postgres "postgres://postgres:postgres@localhost:5432/gorssagg?sslmode=disable" down

.PHONY: migrate-dev-down-all
migrate-dev-down-all:
	goose -dir ./sql/schema postgres "postgres://postgres:postgres@localhost:5432/gorssagg?sslmode=disable" down-to 0

.PHONY: generate-queries
generate-queries:
	sqlc generate

THIS_FILE := $(lastword $(MAKEFILE_LIST))
.PHONY: help build up up-dev start down destroy stop restart logs logs-goserver logs-postgres logs-redis ps login-goserver db-shell login-redis prune
help:
	make -pRrq  -f $(THIS_FILE) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
build:
	docker-compose $(dockerComposeFile) build $(c)
build-dev:
	make build dockerComposeFile="-f docker-compose.dev.yml"
up:
	docker-compose $(dockerComposeFile) up -d $(c)
up-dev:
	docker-compose -f docker-compose.dev.yml up $(c)
start:
	docker-compose start $(c)
down:
	docker-compose down $(c)
destroy:
	docker-compose down -v $(c)
stop:
	docker-compose stop $(c)
restart:
	docker-compose stop $(c)
	docker-compose up -d $(c)
logs:
	docker-compose logs --tail=100 -f $(c)
logs-goserver:
	make logs c=go-app
logs-postgres:
	make logs c=postgres
logs-redis:
	make logs c=redis
ps:
	docker-compose ps
login-goserver:
	docker-compose exec go-app /bin/bash
db-shell:
	docker-compose exec postgres psql -U postgres
login-redis:
	docker-compose exec redis /bin/bash
prune:
	docker image prune -f
	docker container prune -f
	docker builder prune -f
