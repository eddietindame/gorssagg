.PHONY: dev
dev:
	go build -o ./tmp/main ./cmd/main.go && air --build.cmd "go build -o ./tmp/main ./cmd/main.go" --build.bin "./tmp/main"

.PHONY: templ-watch
templ-watch:
	templ generate --watch

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: build
build:
	make templ-generate
	go build -ldflags "-X main.environment=production" -o ./bin/$(APP_NAME) ./cmd/main.go
