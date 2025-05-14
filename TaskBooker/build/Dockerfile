FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && \
    go install github.com/swaggo/swag/cmd/swag@v1.16.2

COPY . .

RUN swag init -g cmd/main.go -o ./docs --parseDependency --parseInternal

RUN CGO_ENABLED=0 go build -o /booker ./cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /booker /booker
COPY --from=builder /app/docs /app/docs
COPY /config/config.yml /app/config/config.yml


EXPOSE 8000

ENTRYPOINT ["/booker"]
