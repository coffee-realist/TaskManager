FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && \
    go install github.com/swaggo/swag/cmd/swag@v1.16.2

COPY . .

RUN swag init -g cmd/main.go -o ./docs --parseDependency --parseInternal

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.2 && \
    rm -rf /var/lib/apt/lists/* && \
    go build -o /publisher ./cmd/main.go

FROM debian:12.10-slim
# FROM scratch
WORKDIR /app

RUN apt-get update && apt-get install -y postgresql-client && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /publisher /publisher
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/docs /app/docs
COPY /config/config.yml /app/config/config.yml

EXPOSE 8080

CMD ["sh", "-c", "migrate -path /app/migrations -database 'postgres://db:db@db:5432/db?sslmode=disable' up && /publisher"]
