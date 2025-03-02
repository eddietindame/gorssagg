FROM golang:1.23.6

WORKDIR /usr/app

COPY . .

RUN apt-get update && apt-get install -y postgresql-client
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go mod download

ENTRYPOINT ["/usr/app/entrypoint.sh"]

CMD ["make", "build-app"]
