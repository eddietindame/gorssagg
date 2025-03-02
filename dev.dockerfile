FROM golang:1.23.6-alpine

WORKDIR /usr/app

COPY . .

RUN chmod +x tailwindcss-linux-arm64-musl
RUN mv tailwindcss-linux-arm64-musl /usr/bin/tailwindcss

RUN apk update
RUN apk upgrade
RUN apk add --no-cache make bash gcc g++ postgresql-client

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/air-verse/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go mod download

ENTRYPOINT ["/usr/app/entrypoint.sh"]

CMD ["make", "dev", "tailwindCmd=tailwindcss"]
