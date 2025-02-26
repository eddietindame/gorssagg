.PHONY: dev
dev:
	go build -o ./tmp/main ./cmd/main.go && air

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: build-server
build-server:
	make templ-generate
	go build -ldflags "-X internal/env.Environment=production" -o ./bin/$(APP_NAME) ./cmd/main.go

THIS_FILE := $(lastword $(MAKEFILE_LIST))
.PHONY: help build up up-dev start down destroy stop restart logs logs-goserver logs-postgres logs-redis ps login-goserver db-shell login-redis
help:
	make -pRrq  -f $(THIS_FILE) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
build:
	docker-compose build $(c)
build-dev:
	docker-compose -f docker-compose.dev.yml build $(c)
up:
	docker-compose up -d $(c)
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
	docker-compose logs --tail=100 -f go-app
logs-postgres:
	docker-compose logs --tail=100 -f postgres
logs-redis:
	docker-compose logs --tail=100 -f redis
ps:
	docker-compose ps
login-goserver:
	docker-compose exec go-app /bin/bash
db-shell:
	docker-compose exec postgres psql -U postgres
login-redis:
	docker-compose exec redis /bin/bash
