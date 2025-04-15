FROM golang:1.24-alpine3.21 AS builder

RUN apk add --no-cache git postgresql-client

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install -tags 'postgres' github.com/golang-migrate/v4/cmd/migrate@v4.18.2 && \
    rm -rf /var/lib/apt/lists/* && \
    go build -o /publisher ./cmd/main.go

FROM debian:12.10-slim
# FROM scratch
RUN apt-get update && apt-get install -y postgresql-client && \
    rm -rf /var/lib/apt/lists/* \

COPY --from=builder /publisher /publisher
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8000

CMD ["migrate -path /app/migrations -database 'postgres://db:db@db:5432/postgres?sslmode=disable' up && /publisher"]
