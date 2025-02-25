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
